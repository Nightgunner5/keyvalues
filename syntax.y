%{
package keyvalues
%}

%union {
	kv *KeyValues
	s  string
	m  map[string]interface{}
}

%type <kv> keyvalues key
%type <m> keys
%type <s>  val strs
%token <s> str

%left '\"' '{' '}'

%start keyvalues

%%

keyvalues: key
	{ $$, yylex.(*parser).kv = $1, $1 }
;

key:
	val '{' keys '}'
		{ $$ = &KeyValues{name: $1, data: $3} }
;

val:
	'\"' strs '\"'
		{ $$ = $2 }
|	str
		{ $$ = $1 }

strs:
	str
		{ $$ = $1 }
|	strs str
		{ $$ = $1 + " " + $2 }
;

keys:
	keys key
		{ $$ = $1; $$[$2.name] = $2 }
|	keys val val
		{ $$ = $1; $$[$2] = $3 }
|	/* empty */
		{ $$ = make(map[string]interface{}) }
;
