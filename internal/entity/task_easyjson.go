// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package entity

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

func easyjson79a0a577DecodeGithubComK1nkyTookhookInternalEntity(in *jlexer.Lexer, out *QueueTask) {
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
		case "Queue":
			out.Queue = string(in.String())
		case "Payload":
			if in.IsNull() {
				in.Skip()
				out.Payload = nil
			} else {
				out.Payload = in.Bytes()
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
func easyjson79a0a577EncodeGithubComK1nkyTookhookInternalEntity(out *jwriter.Writer, in QueueTask) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Queue\":"
		out.RawString(prefix[1:])
		out.String(string(in.Queue))
	}
	{
		const prefix string = ",\"Payload\":"
		out.RawString(prefix)
		out.Base64Bytes(in.Payload)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v QueueTask) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson79a0a577EncodeGithubComK1nkyTookhookInternalEntity(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v QueueTask) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson79a0a577EncodeGithubComK1nkyTookhookInternalEntity(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *QueueTask) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson79a0a577DecodeGithubComK1nkyTookhookInternalEntity(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *QueueTask) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson79a0a577DecodeGithubComK1nkyTookhookInternalEntity(l, v)
}
