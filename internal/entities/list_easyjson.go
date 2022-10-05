// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package entities

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

func easyjsonEf7cfe30DecodeGolangTestTaskInternalEntities(in *jlexer.Lexer, out *ListAdsAnswer) {
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
		case "status":
			out.Status = string(in.String())
		case "result":
			if in.IsNull() {
				in.Skip()
				out.Result = nil
			} else {
				in.Delim('[')
				if out.Result == nil {
					if !in.IsDelim(']') {
						out.Result = make([]APIAdListItem, 0, 1)
					} else {
						out.Result = []APIAdListItem{}
					}
				} else {
					out.Result = (out.Result)[:0]
				}
				for !in.IsDelim(']') {
					var v1 APIAdListItem
					(v1).UnmarshalEasyJSON(in)
					out.Result = append(out.Result, v1)
					in.WantComma()
				}
				in.Delim(']')
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
func easyjsonEf7cfe30EncodeGolangTestTaskInternalEntities(out *jwriter.Writer, in ListAdsAnswer) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"status\":"
		out.RawString(prefix[1:])
		out.String(string(in.Status))
	}
	{
		const prefix string = ",\"result\":"
		out.RawString(prefix)
		if in.Result == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Result {
				if v2 > 0 {
					out.RawByte(',')
				}
				(v3).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ListAdsAnswer) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonEf7cfe30EncodeGolangTestTaskInternalEntities(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ListAdsAnswer) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonEf7cfe30EncodeGolangTestTaskInternalEntities(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ListAdsAnswer) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonEf7cfe30DecodeGolangTestTaskInternalEntities(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ListAdsAnswer) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonEf7cfe30DecodeGolangTestTaskInternalEntities(l, v)
}
func easyjsonEf7cfe30DecodeGolangTestTaskInternalEntities1(in *jlexer.Lexer, out *APIAdListItem) {
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
			out.ID = int(in.Int())
		case "title":
			out.Title = string(in.String())
		case "price":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Price).UnmarshalJSON(data))
			}
		case "main_image_url":
			if in.IsNull() {
				in.Skip()
				out.MainImageURL = nil
			} else {
				if out.MainImageURL == nil {
					out.MainImageURL = new(string)
				}
				*out.MainImageURL = string(in.String())
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
func easyjsonEf7cfe30EncodeGolangTestTaskInternalEntities1(out *jwriter.Writer, in APIAdListItem) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int(int(in.ID))
	}
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix)
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"price\":"
		out.RawString(prefix)
		out.Raw((in.Price).MarshalJSON())
	}
	{
		const prefix string = ",\"main_image_url\":"
		out.RawString(prefix)
		if in.MainImageURL == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.MainImageURL))
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v APIAdListItem) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonEf7cfe30EncodeGolangTestTaskInternalEntities1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v APIAdListItem) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonEf7cfe30EncodeGolangTestTaskInternalEntities1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *APIAdListItem) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonEf7cfe30DecodeGolangTestTaskInternalEntities1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *APIAdListItem) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonEf7cfe30DecodeGolangTestTaskInternalEntities1(l, v)
}
