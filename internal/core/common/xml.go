package common

import (
	"strings"
)

// XmlBuilder is a simple utility to build XML strings.
type XmlBuilder struct {
	sb strings.Builder
	stack []string
}

func NewXmlBuilder() *XmlBuilder {
	return &XmlBuilder{
		stack: []string{},
	}
}

func (b *XmlBuilder) Start(name string) *XmlBuilder {
	b.sb.WriteString("<")
	b.sb.WriteString(name)
	b.sb.WriteString(">")
	b.stack = append(b.stack, name)
	return b
}

func (b *XmlBuilder) End() *XmlBuilder {
	if len(b.stack) == 0 {
		return b
	}
	name := b.stack[len(b.stack)-1]
	b.stack = b.stack[:len(b.stack)-1]
	b.sb.WriteString("</")
	b.sb.WriteString(name)
	b.sb.WriteString(">")
	return b
}

func (b *XmlBuilder) Elem(name, value string) *XmlBuilder {
	b.sb.WriteString("<")
	b.sb.WriteString(name)
	b.sb.WriteString(">")
	b.sb.WriteString(value) // TODO: XML escape
	b.sb.WriteString("</")
	b.sb.WriteString(name)
	b.sb.WriteString(">")
	return b
}

func (b *XmlBuilder) Raw(xml string) *XmlBuilder {
	b.sb.WriteString(xml)
	return b
}


func (b *XmlBuilder) Build() string {
	for len(b.stack) > 0 {
		b.End()
	}
	return b.sb.String()
}
