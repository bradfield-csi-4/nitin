```
query     : or_query

or_query  : and_query ('OR' and_query)*

and_query : not_query ('AND'? not_query)*

not_query : ('NOT'|'-')* atom

atom      : term
          | phrase
          | '(' query ')'
```
