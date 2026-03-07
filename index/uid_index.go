package index

import (
	"context"
	"iter"
	"sort"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"go.linka.cloud/protofilters/filters"
	"go.linka.cloud/protofilters/index/bitmap"
	preflect "go.linka.cloud/protofilters/reflect"
)

type FindOptions struct {
	Offset  uint64
	Limit   uint64
	Reverse bool
}

// UIDIndex is a protobuf message index keyed by UID.
type UIDIndex interface {
	Insert(ctx context.Context, uid uint64, m proto.Message) error
	Update(ctx context.Context, uid uint64, old, m proto.Message) error
	Remove(ctx context.Context, uid uint64) error
	Find(ctx context.Context, t protoreflect.FullName, f filters.FieldFilterer, opts FindOptions) iter.Seq2[uint64, error]
}

type uidIndex struct {
	store UIDTxer
	fn    Func
}

// NewUID creates a new UID index using the given store and index function.
func NewUID(s UIDStore, fn Func) UIDIndex {
	if fn == nil {
		fn = All
	}
	if s == nil {
		s = newUIDStore()
	}
	x, ok := any(s).(UIDTxer)
	if !ok {
		x = &fakeUIDTxer{UIDStore: s}
	}
	return &uidIndex{store: x, fn: fn}
}

func newUIDFromTxer(x UIDTxer, fn Func) UIDIndex {
	if fn == nil {
		fn = All
	}
	return &uidIndex{store: x, fn: fn}
}

func (i *uidIndex) index(ctx context.Context, tx UIDTx, uid uint64, m protoreflect.Message, fds ...protoreflect.FieldDescriptor) error {
	f := m.Descriptor().Fields()
	name := m.Descriptor().FullName()
	for j := 0; j < f.Len(); j++ {
		fd := f.Get(j)
		path := append(fds, fd)
		ok, err := i.fn(ctx, name, path...)
		if err != nil {
			return err
		}
		if isUnsetRealOneofField(m, fd) {
			continue
		}
		rval := m.Get(fd)
		if fd.IsList() {
			if fd.Kind() == protoreflect.MessageKind {
				for j2 := 0; j2 < rval.List().Len(); j2++ {
					if err := i.index(ctx, tx, uid, rval.List().Get(j2).Message(), path...); err != nil {
						return err
					}
				}
				continue
			}
			list := rval.List()
			for j2 := 0; j2 < list.Len(); j2++ {
				if err := tx.AddUID(ctx, uid, list.Get(j2), path...); err != nil {
					return err
				}
			}
			continue
		}
		if fd.IsMap() {
			continue
		}
		if fd.Kind() == protoreflect.MessageKind && !preflect.IsWKType(fd.Message().FullName()) {
			if !rval.Message().IsValid() {
				continue
			}
			if err := i.index(ctx, tx, uid, rval.Message(), path...); err != nil {
				return err
			}
			continue
		}
		if fd.HasOptionalKeyword() && !m.Has(fd) {
			rval = protoreflect.Value{}
		}
		if !ok {
			continue
		}
		if err := tx.AddUID(ctx, uid, rval, path...); err != nil {
			return err
		}
	}
	return nil
}

func (i *uidIndex) Insert(ctx context.Context, uid uint64, m proto.Message) error {
	tx, err := i.store.Tx(ctx)
	if err != nil {
		return err
	}
	defer tx.Close()
	if err := i.index(ctx, tx, uid, m.ProtoReflect()); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (i *uidIndex) Update(ctx context.Context, uid uint64, old, m proto.Message) error {
	tx, err := i.store.Tx(ctx)
	if err != nil {
		return err
	}
	defer tx.Close()
	oldValues := map[string]fieldValues{}
	newValues := map[string]fieldValues{}
	if old != nil {
		if oldValues, err = i.collectValues(ctx, old.ProtoReflect()); err != nil {
			return err
		}
	}
	if m != nil {
		if newValues, err = i.collectValues(ctx, m.ProtoReflect()); err != nil {
			return err
		}
	}
	if len(oldValues) > 0 {
		fr, err := tx.For(ctx, old.ProtoReflect().Descriptor().FullName())
		if err != nil {
			return err
		}
		if err := applyUIDDiff(ctx, fr, tx, uid, oldValues, newValues); err != nil {
			return err
		}
	} else {
		for _, fv := range newValues {
			for _, v := range fv.values {
				if err := tx.AddUID(ctx, uid, v, fv.fds...); err != nil {
					return err
				}
			}
		}
	}
	return tx.Commit(ctx)
}

func (i *uidIndex) collectValues(ctx context.Context, m protoreflect.Message) (map[string]fieldValues, error) {
	values := map[string]fieldValues{}
	if err := i.collectValuesInto(ctx, values, m); err != nil {
		return nil, err
	}
	return values, nil
}

func (i *uidIndex) collectValuesInto(ctx context.Context, out map[string]fieldValues, m protoreflect.Message, fds ...protoreflect.FieldDescriptor) error {
	f := m.Descriptor().Fields()
	name := m.Descriptor().FullName()
	for j := 0; j < f.Len(); j++ {
		fd := f.Get(j)
		path := append(fds, fd)
		ok, err := i.fn(ctx, name, path...)
		if err != nil {
			return err
		}
		if isUnsetRealOneofField(m, fd) {
			continue
		}
		rval := m.Get(fd)
		if fd.IsList() {
			if fd.Kind() == protoreflect.MessageKind {
				for j2 := 0; j2 < rval.List().Len(); j2++ {
					if err := i.collectValuesInto(ctx, out, rval.List().Get(j2).Message(), path...); err != nil {
						return err
					}
				}
				continue
			}
			list := rval.List()
			for j2 := 0; j2 < list.Len(); j2++ {
				out = appendValue(out, path, list.Get(j2))
			}
			continue
		}
		if fd.IsMap() {
			continue
		}
		if fd.Kind() == protoreflect.MessageKind && !preflect.IsWKType(fd.Message().FullName()) {
			if !rval.Message().IsValid() {
				continue
			}
			if err := i.collectValuesInto(ctx, out, rval.Message(), path...); err != nil {
				return err
			}
			continue
		}
		if fd.HasOptionalKeyword() && !m.Has(fd) {
			rval = protoreflect.Value{}
		}
		if !ok {
			continue
		}
		out = appendValue(out, path, rval)
	}
	return nil
}

func applyUIDDiff(ctx context.Context, fr FieldReader, tx UIDTx, uid uint64, oldValues, newValues map[string]fieldValues) error {
	seen := map[string]struct{}{}
	for key := range oldValues {
		seen[key] = struct{}{}
	}
	for key := range newValues {
		seen[key] = struct{}{}
	}
	for key := range seen {
		ov := oldValues[key]
		nv := newValues[key]
		remove, add := diffValues(ov.values, nv.values)
		for _, v := range remove {
			if err := removeUIDValue(ctx, fr, uid, v, ov.fds...); err != nil {
				return err
			}
		}
		for _, v := range add {
			if err := tx.AddUID(ctx, uid, v, nv.fds...); err != nil {
				return err
			}
		}
	}
	return nil
}

func removeUIDValue(ctx context.Context, fr FieldReader, uid uint64, v protoreflect.Value, fds ...protoreflect.FieldDescriptor) error {
	if len(fds) == 0 {
		return nil
	}
	name := joinFieldNames(fds)
	for f, err := range fr.Get(ctx, name) {
		if err != nil {
			return err
		}
		if !valueEqual(f.Value(), v) {
			continue
		}
		b, err := f.Bitmap(ctx)
		if err != nil {
			return err
		}
		b.Remove(uid)
		return nil
	}
	return nil
}

func (i *uidIndex) Remove(ctx context.Context, uid uint64) error {
	tx, err := i.store.Tx(ctx)
	if err != nil {
		return err
	}
	defer tx.Close()
	if err := tx.ClearUID(ctx, uid); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (i *uidIndex) doFind(ctx context.Context, tx UIDTx, t protoreflect.FullName, f *filters.FieldFilter) (bitmap.Bitmap, error) {
	fds, err := tx.For(ctx, t)
	if err != nil {
		return nil, err
	}
	b := bitmap.NewWith(1024)
	for v, err := range fds.Get(ctx, protoreflect.Name(f.Field)) {
		if err != nil {
			return nil, err
		}
		ds := v.Descriptors()
		fd := ds[len(ds)-1]
		ok, err := preflect.Match(v.Value(), fd, f.Filter)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}
		b2, err := v.Bitmap(ctx)
		if err != nil {
			return nil, err
		}
		b.Or(b2)
	}
	return b, nil
}

func (i *uidIndex) find(ctx context.Context, tx UIDTx, t protoreflect.FullName, f filters.FieldFilterer) (bitmap.Bitmap, error) {
	expr := f.Expr()
	b, err := i.doFind(ctx, tx, t, expr.Condition)
	if err != nil {
		return nil, err
	}
	for _, v := range expr.AndExprs {
		b2, err := i.find(ctx, tx, t, v)
		if err != nil {
			return nil, err
		}
		b.And(b2)
	}
	for _, v := range expr.OrExprs {
		b2, err := i.find(ctx, tx, t, v)
		if err != nil {
			return nil, err
		}
		b.Or(b2)
	}
	return b, nil
}

func (i *uidIndex) Find(ctx context.Context, t protoreflect.FullName, f filters.FieldFilterer, opts FindOptions) iter.Seq2[uint64, error] {
	return func(yield func(uint64, error) bool) {
		if f == nil || f.Expr() == nil {
			return
		}
		tx, err := i.store.Tx(ctx)
		if err != nil {
			yield(0, err)
			return
		}
		defer tx.Close()
		b, err := i.find(ctx, tx, t, f)
		if err != nil {
			yield(0, err)
			return
		}

		uids := make([]uint64, 0, b.Cardinality())
		for uid := range b.Iter() {
			uids = append(uids, uid)
		}
		sort.Slice(uids, func(a, b int) bool { return uids[a] < uids[b] })

		var emitted uint64
		if !opts.Reverse {
			var skipped uint64
			for _, uid := range uids {
				if skipped < opts.Offset {
					skipped++
					continue
				}
				if opts.Limit > 0 && emitted >= opts.Limit {
					return
				}
				emitted++
				if !yield(uid, nil) {
					return
				}
			}
			return
		}

		var skipped uint64
		for idx := len(uids) - 1; idx >= 0; idx-- {
			if skipped < opts.Offset {
				skipped++
				continue
			}
			if opts.Limit > 0 && emitted >= opts.Limit {
				return
			}
			emitted++
			if !yield(uids[idx], nil) {
				return
			}
		}
	}
}
