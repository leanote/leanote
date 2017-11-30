# Ace 
 The [ACE](https://github.com/yosssi/ace) Templating Plugin
 
- Ace templates have full access to the revel.TemplateFuncs, any function 
 defined in there can be used by the this engine

- Ace files must be identified by using a `shebang` on the first line 
(preferred method) or changing the file extension to home.ace.html. 
 By default the operation of revel assumes just a `.html` extension for
 controller responses so it would make more sense to use the shebang
- Ace templates can be set to be case sensitive by setting
`ace.tempate.caseinsensitive=false`, default is not case sensitive. If case sensitivity
is off internal imports must be done using lower case
- All function registered in `revel.TemplateFuncs` are available for use 
inside the ace framework

##### Details
Ace is a little different of a templating system, its output is a 
standard go template but there is no concept of template sets, 
instead you build a composite template using
 a *base* template and an *inner* template. The 
 *inner* template can only contain items like : 
   ```
= content main
  h2 Inner Template - Main : {{.Msg}}

= content sub
  h3 Inner Template - Sub : {{.Msg}}
     
   ```
The base template can contain items like 
```
= doctype html
html lang=en
  head
    meta charset=utf-8
    title Ace example
    = css
      h1 { color: blue; }
  body
    h1 Base Template : {{.Msg}}
    #container.wrapper
      = yield main
      = yield sub
      = include inc .Msg
    = javascript
      alert('{{.Msg}}');
```

You are allowed to include one *inner* template with the base template,
to do so in revel you can extend your controller from the ace controller
and call `RenderAceTemplate(base ,inner string)` which will insert
the inner template using the outer template.
 
 The ace engine requires that you explicitly set the template type on the
 template itself by either using the shebang method on the first line
 like `#! ace` or having the file name like `template.ace.html` 
 either method will work. 

