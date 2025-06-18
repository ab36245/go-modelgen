package dartgen

import (
	"github.com/ab36245/go-modelgen/defx"
)

func newModels(ds []defx.Model, opts Opts) *modelsGen {
	list := []*modelGen{}
	for _, d := range ds {
		list = append(list, newModel(d))
	}
	return &modelsGen{
		list: list,
		opts: opts,
	}
}

type modelsGen struct {
	list []*modelGen
	opts Opts
}

func (g *modelsGen) code() string {
	put("// WARNING!")
	put("// This code was generated automatically.")
	put("")
	put("import 'package:flutter_model/flutter_model.dart';")
	put("import 'package:flutter_writer/flutter_writer.dart';")

	for _, m := range g.list {
		put("")
		m.doClass()
	}

	// put("")
	// inc("Msg _decoders(int id, MsgPackDecoder mp) =>")
	// {
	// 	inc("switch (id) {")
	// 	{
	// 		for _, m := range g.list {
	// 			put("%s => _decode%s(mp),", m.id, m.name)
	// 		}
	// 		put("_ => throw MsgsException('unknown msg id $id')")
	// 	}
	// 	dec("};")
	// }
	// dec("")

	// for _, m := range g.list {
	// 	put("")
	// 	m.doCodec()
	// }

	return w.String()
}
