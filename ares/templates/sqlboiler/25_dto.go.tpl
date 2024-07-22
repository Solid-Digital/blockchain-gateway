{{- $alias := .Aliases.Table .Table.Name -}}

// {{$alias.UpSingular}} is an object representing the database table.
type {{$alias.UpSingular}}DTO struct {
	{{- range $column := .Table.Columns -}}
	{{- $colAlias := $alias.Column $column.Name -}}
	{{- if eq $.StructTagCasing "camel"}}
	{{$colAlias}} {{$column.Type}} `{{generateTags $.Tags $column.Name}}boil:"{{$column.Name}}" json:"{{$column.Name | camelCase}}{{if $column.Nullable}},omitempty{{end}}" toml:"{{$column.Name | camelCase}}" yaml:"{{$column.Name | camelCase}}{{if $column.Nullable}},omitempty{{end}}"`
	{{- else -}}
	{{$colAlias}} {{$column.Type}} `{{generateTags $.Tags $column.Name}}boil:"{{$column.Name}}" json:"{{$column.Name}}{{if $column.Nullable}},omitempty{{end}}" toml:"{{$column.Name}}" yaml:"{{$column.Name}}{{if $column.Nullable}},omitempty{{end}}"`
	{{end -}}
	{{end -}}
}

// DTO converts the {{$alias.UpSingular}} to a {{$alias.UpSingular}}DTO struct.
func (o {{$alias.UpSingular}}) DTO() (*{{$alias.UpSingular}}DTO) {
    return &{{$alias.UpSingular}}DTO{
    {{ range $column := .Table.Columns -}}
    {{- $colAlias := $alias.Column $column.Name -}}
    {{ $colAlias}}: o.{{$colAlias}},
    {{end -}}
    }
}
