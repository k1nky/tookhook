// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package builtin

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

func easyjsonB5edef53DecodeGithubComK1nkyTookhookInternalBuiltin(in *jlexer.Lexer, out *ExecHandlerOptions) {
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
		case "shell":
			out.Shell = string(in.String())
		case "args":
			if in.IsNull() {
				in.Skip()
				out.Args = nil
			} else {
				in.Delim('[')
				if out.Args == nil {
					if !in.IsDelim(']') {
						out.Args = make([]string, 0, 4)
					} else {
						out.Args = []string{}
					}
				} else {
					out.Args = (out.Args)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.Args = append(out.Args, v1)
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
func easyjsonB5edef53EncodeGithubComK1nkyTookhookInternalBuiltin(out *jwriter.Writer, in ExecHandlerOptions) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"shell\":"
		out.RawString(prefix[1:])
		out.String(string(in.Shell))
	}
	{
		const prefix string = ",\"args\":"
		out.RawString(prefix)
		if in.Args == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Args {
				if v2 > 0 {
					out.RawByte(',')
				}
				out.String(string(v3))
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ExecHandlerOptions) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonB5edef53EncodeGithubComK1nkyTookhookInternalBuiltin(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ExecHandlerOptions) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonB5edef53EncodeGithubComK1nkyTookhookInternalBuiltin(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ExecHandlerOptions) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonB5edef53DecodeGithubComK1nkyTookhookInternalBuiltin(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ExecHandlerOptions) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonB5edef53DecodeGithubComK1nkyTookhookInternalBuiltin(l, v)
}
