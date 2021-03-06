PVNT_GO
=========

### Introduction

_PVNT is a simple, open source template utility written in Go Language. It is GOLang port for [PVNT](https://github.com/pvnhome/pvnt) written in Java._

* Version: 0.94 (beta)
* Author: [pvnhome](mailto:pvnhome@yandex.ru)

### Features

In very simple case we need one template file and several HTML-files that use this template (for example A.html and B.html).

Template file (in _templates_ subdirectory):
```html
<!DOCTYPE html>
<html lang="ru">
<!--pvnTmplBeg-->
<head>
   <meta charset="utf-8">
   <!--pvnEditBeg title-->
   <title>PageTemplate</title>
   <!--pvnEditEnd title-->
</head>
<body>
   <h1>Site1 (A,B)</h1>
   <h2>Part1</h2>
   <!--pvnEditBeg part1-->
   <p>Palceholder for Part1</p>
   <!--pvnEditEnd part1-->
   <h2>Part2</h2>
   <!--pvnEditBeg part2-->
   <p>Palceholder for Part2</p>
   <!--pvnEditEnd part2-->
</body>
<!--pvnTmplEnd-->
</html>
```

HTML-file A.html:
```html
<!DOCTYPE html>
<html lang="ru">
<!--pvnTmplBeg templates/template.html-->
<head>
   <meta charset="utf-8">
   <!--pvnImplBeg title-->
   <title>PageA</title>
   <!--pvnImplEnd title-->
</head>
<body>
   <h1>Site1 (A,B)</h1>
   <h2>Part1</h2>
   <!--pvnImplBeg part1-->
   <p>Content for Part1 on PageA.</p>
   <!--pvnImplEnd part1-->
   <h2>Part2</h2>
   <!--pvnImplBeg part2-->
   <p>Content for Part2 on PageA.</p>
   <!--pvnImplEnd part2-->
</body>
<!--pvnTmplEnd-->
</html>
```
HTML-file B.html:
```html
<!DOCTYPE html>
<html lang="ru">
<!--pvnTmplBeg templates/template.html-->
<head>
   <meta charset="utf-8">
   <!--pvnImplBeg title-->
   <title>PageB</title>
   <!--pvnImplEnd title-->
</head>
<body>
   <h1>Site1 (A,B)</h1>
   <h2>Part1</h2>
   <!--pvnImplBeg part1-->
   <p>Content for Part1 on PageB.</p>
   <!--pvnImplEnd part1-->
   <h2>Part2</h2>
   <!--pvnImplBeg part2-->
   <p>Content for Part2 on PageB.</p>
   <!--pvnImplEnd part2-->
</body>
<!--pvnTmplEnd-->
</html>
```
We can use only three types of tags:
* pvnTmpl - used to define template file
* pvnEdit - used to define editable region in template file
* pvnImpl - used to override editable region

### Installation

Installation as simple as 1-2-3:
* Download _pvnt_go-x.x.x.zip_ or _pvnt_go-x.x.x-linux-amd64.tar.gz_ and extract it into the directory of your choice.
* Add the _[PVNT]/bin_ directory to your _PATH_ environment variable.
* Change _[PVNT]/examples/site/templates/template.html_ file.
* Execute `pvnt_go` command in _[PVNT]/examples/site_ directory.
* Look at A.html and B.html files for changes.

### Website
* PVNT in Go Language Source code
<https://github.com/pvnhome/pvnt_go>
* PVNT in Java Source code
<https://github.com/pvnhome/pvnt>
* Binary downloads 
<https://github.com/pvnhome/pvnt_go/releases/tag/v0.9-beta.4>
