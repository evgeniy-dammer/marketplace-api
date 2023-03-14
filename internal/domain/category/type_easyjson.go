// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package category

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonBc289ab0DecodeGithubComEvgeniyDammermarketplaceApiInternalDomainCategory(in *jlexer.Lexer, out *UpdateCategoryInput) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			if in.IsNull() {
				in.Skip()
				out.ID = nil
			} else {
				if out.ID == nil {
					out.ID = new(string)
				}
				*out.ID = string(in.String())
			}
		case "nametm":
			if in.IsNull() {
				in.Skip()
				out.NameTm = nil
			} else {
				if out.NameTm == nil {
					out.NameTm = new(string)
				}
				*out.NameTm = string(in.String())
			}
		case "nameru":
			if in.IsNull() {
				in.Skip()
				out.NameRu = nil
			} else {
				if out.NameRu == nil {
					out.NameRu = new(string)
				}
				*out.NameRu = string(in.String())
			}
		case "nametr":
			if in.IsNull() {
				in.Skip()
				out.NameTr = nil
			} else {
				if out.NameTr == nil {
					out.NameTr = new(string)
				}
				*out.NameTr = string(in.String())
			}
		case "nameen":
			if in.IsNull() {
				in.Skip()
				out.NameEn = nil
			} else {
				if out.NameEn == nil {
					out.NameEn = new(string)
				}
				*out.NameEn = string(in.String())
			}
		case "parent":
			if in.IsNull() {
				in.Skip()
				out.Parent = nil
			} else {
				if out.Parent == nil {
					out.Parent = new(string)
				}
				*out.Parent = string(in.String())
			}
		case "organisation":
			if in.IsNull() {
				in.Skip()
				out.OrganizationID = nil
			} else {
				if out.OrganizationID == nil {
					out.OrganizationID = new(string)
				}
				*out.OrganizationID = string(in.String())
			}
		case "level":
			if in.IsNull() {
				in.Skip()
				out.Level = nil
			} else {
				if out.Level == nil {
					out.Level = new(int)
				}
				*out.Level = int(in.Int())
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonBc289ab0EncodeGithubComEvgeniyDammermarketplaceApiInternalDomainCategory(out *jwriter.Writer, in UpdateCategoryInput) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		if in.ID == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.ID))
		}
	}
	{
		const prefix string = ",\"nametm\":"
		out.RawString(prefix)
		if in.NameTm == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.NameTm))
		}
	}
	{
		const prefix string = ",\"nameru\":"
		out.RawString(prefix)
		if in.NameRu == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.NameRu))
		}
	}
	{
		const prefix string = ",\"nametr\":"
		out.RawString(prefix)
		if in.NameTr == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.NameTr))
		}
	}
	{
		const prefix string = ",\"nameen\":"
		out.RawString(prefix)
		if in.NameEn == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.NameEn))
		}
	}
	{
		const prefix string = ",\"parent\":"
		out.RawString(prefix)
		if in.Parent == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.Parent))
		}
	}
	{
		const prefix string = ",\"organisation\":"
		out.RawString(prefix)
		if in.OrganizationID == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.OrganizationID))
		}
	}
	{
		const prefix string = ",\"level\":"
		out.RawString(prefix)
		if in.Level == nil {
			out.RawString("null")
		} else {
			out.Int(int(*in.Level))
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UpdateCategoryInput) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonBc289ab0EncodeGithubComEvgeniyDammermarketplaceApiInternalDomainCategory(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UpdateCategoryInput) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonBc289ab0EncodeGithubComEvgeniyDammermarketplaceApiInternalDomainCategory(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UpdateCategoryInput) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonBc289ab0DecodeGithubComEvgeniyDammermarketplaceApiInternalDomainCategory(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UpdateCategoryInput) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonBc289ab0DecodeGithubComEvgeniyDammermarketplaceApiInternalDomainCategory(l, v)
}
func easyjsonBc289ab0DecodeGithubComEvgeniyDammermarketplaceApiInternalDomainCategory1(in *jlexer.Lexer, out *ListCategory) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(ListCategory, 0, 0)
			} else {
				*out = ListCategory{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 Category
			(v1).UnmarshalEasyJSON(in)
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonBc289ab0EncodeGithubComEvgeniyDammermarketplaceApiInternalDomainCategory1(out *jwriter.Writer, in ListCategory) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			(v3).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v ListCategory) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonBc289ab0EncodeGithubComEvgeniyDammermarketplaceApiInternalDomainCategory1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ListCategory) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonBc289ab0EncodeGithubComEvgeniyDammermarketplaceApiInternalDomainCategory1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ListCategory) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonBc289ab0DecodeGithubComEvgeniyDammermarketplaceApiInternalDomainCategory1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ListCategory) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonBc289ab0DecodeGithubComEvgeniyDammermarketplaceApiInternalDomainCategory1(l, v)
}
func easyjsonBc289ab0DecodeGithubComEvgeniyDammermarketplaceApiInternalDomainCategory2(in *jlexer.Lexer, out *CreateCategoryInput) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "nametm":
			out.NameTm = string(in.String())
		case "nameru":
			out.NameRu = string(in.String())
		case "nametr":
			out.NameTr = string(in.String())
		case "nameen":
			out.NameEn = string(in.String())
		case "parent":
			out.Parent = string(in.String())
		case "organization":
			out.OrganizationID = string(in.String())
		case "level":
			out.Level = int(in.Int())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonBc289ab0EncodeGithubComEvgeniyDammermarketplaceApiInternalDomainCategory2(out *jwriter.Writer, in CreateCategoryInput) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"nametm\":"
		out.RawString(prefix[1:])
		out.String(string(in.NameTm))
	}
	{
		const prefix string = ",\"nameru\":"
		out.RawString(prefix)
		out.String(string(in.NameRu))
	}
	{
		const prefix string = ",\"nametr\":"
		out.RawString(prefix)
		out.String(string(in.NameTr))
	}
	{
		const prefix string = ",\"nameen\":"
		out.RawString(prefix)
		out.String(string(in.NameEn))
	}
	{
		const prefix string = ",\"parent\":"
		out.RawString(prefix)
		out.String(string(in.Parent))
	}
	{
		const prefix string = ",\"organization\":"
		out.RawString(prefix)
		out.String(string(in.OrganizationID))
	}
	{
		const prefix string = ",\"level\":"
		out.RawString(prefix)
		out.Int(int(in.Level))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CreateCategoryInput) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonBc289ab0EncodeGithubComEvgeniyDammermarketplaceApiInternalDomainCategory2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CreateCategoryInput) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonBc289ab0EncodeGithubComEvgeniyDammermarketplaceApiInternalDomainCategory2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CreateCategoryInput) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonBc289ab0DecodeGithubComEvgeniyDammermarketplaceApiInternalDomainCategory2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CreateCategoryInput) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonBc289ab0DecodeGithubComEvgeniyDammermarketplaceApiInternalDomainCategory2(l, v)
}
func easyjsonBc289ab0DecodeGithubComEvgeniyDammermarketplaceApiInternalDomainCategory3(in *jlexer.Lexer, out *Category) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = string(in.String())
		case "nametm":
			out.NameTm = string(in.String())
		case "nameru":
			out.NameRu = string(in.String())
		case "nametr":
			out.NameTr = string(in.String())
		case "nameen":
			out.NameEn = string(in.String())
		case "parent":
			out.Parent = string(in.String())
		case "organization":
			out.OrganizationID = string(in.String())
		case "level":
			out.Level = int(in.Int())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonBc289ab0EncodeGithubComEvgeniyDammermarketplaceApiInternalDomainCategory3(out *jwriter.Writer, in Category) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.String(string(in.ID))
	}
	{
		const prefix string = ",\"nametm\":"
		out.RawString(prefix)
		out.String(string(in.NameTm))
	}
	{
		const prefix string = ",\"nameru\":"
		out.RawString(prefix)
		out.String(string(in.NameRu))
	}
	{
		const prefix string = ",\"nametr\":"
		out.RawString(prefix)
		out.String(string(in.NameTr))
	}
	{
		const prefix string = ",\"nameen\":"
		out.RawString(prefix)
		out.String(string(in.NameEn))
	}
	{
		const prefix string = ",\"parent\":"
		out.RawString(prefix)
		out.String(string(in.Parent))
	}
	{
		const prefix string = ",\"organization\":"
		out.RawString(prefix)
		out.String(string(in.OrganizationID))
	}
	{
		const prefix string = ",\"level\":"
		out.RawString(prefix)
		out.Int(int(in.Level))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Category) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonBc289ab0EncodeGithubComEvgeniyDammermarketplaceApiInternalDomainCategory3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Category) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonBc289ab0EncodeGithubComEvgeniyDammermarketplaceApiInternalDomainCategory3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Category) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonBc289ab0DecodeGithubComEvgeniyDammermarketplaceApiInternalDomainCategory3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Category) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonBc289ab0DecodeGithubComEvgeniyDammermarketplaceApiInternalDomainCategory3(l, v)
}
