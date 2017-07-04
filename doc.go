// Copyright 2017 The Ottemo Authors. All rights reserved.

/*
Package foundation represents the Ottemo e-commerce application entry point. This package contains the go main()
function which assembles the server application.

  Example:
    go build github.com/ottemo/foundation
    go build -tags mongo github.com/ottemo/foundation


Project structure and convections outline:

1. Packages structure:

	    It is well  known that GO packages represented by  directories. Within package you can have  zero or more files,
	which considered by GO as  one build unit with shared scope. There are no strict  naming convection GO lang proposes
	to file names, however there is a convection of package  naming that says you should try to name in single-word (see
	https://golang.org/doc/effective_go.html#names).

	    Despite the fact  that there is no  package files naming convection for  GO, there are couple  limits you should
	consider: file  names ends with  "_test.go" (refer to http://golang.org/pkg/testing/#pkg-overview)  and architecture
	files ends  on "*_GOOS", "*_GOARCH", "*_GOOS_GOARCH"  (i.e.  like some_amd64.go, some_windows_amd64.go,  on possible
	their values you can refer to https://golang.org/doc/install/source#environment, http://golang.org/pkg/go/build/).

	    So, these files are treated different way as well as files prefixed by "." which are considered non existent.

	    Ottemo project  trying to  follow Google recommendations,  but - with  saving clear  on code units.   So, Ottemo
	packages have multi-level folder structure with little convection on GO file naming within package.

	      file name          |  contents / meaning
	    ---------------------|-------------------------------------------------------------------------------
	    decl.go              | declarations of structures/variables/constants package wants to share
	    init.go              | package initialization routines, i.e. init() function
	                         | (refer http://golang.org/ref/spec#Package_initialization)
                             |
	    doc.go               | package description for GO doc
	                         |
	    config.go            | configuration values and their validators declaration
	    interfaces.go        | declaration of interfaces provided by package
	    manager.go           | routines for register/unregister package interface implementations
	                         |
	    helpers.go           | shortcuts for objects creation, repeatable routines, etc.
	    utils.go             | tools to work with package objects, but could be other purposes
	                         |
	    [i_interface]        | set of functions implementing particular interface (one unit case)
	    [unit].[i_interface] | set of functions implementing particular interface (multi unit case)
	                         |
	    [name]               | package specific operational logic separation (one unit case)
	    [unit].[name]        | package specific operational logic separation (multi unit case)

	    [name] and [unit] are not something specific - just some word to reflect you have different scope.

	    Considering these naming  agreements, you can determinate package  role based on just listing  package files and
	examining  of "interfaces.go"  and  "decl.go"  files.  Also  you  can consider  amount  of  interfaces this  package
	implements, by counting "i_" prefixed files.


2. Packages types/roles:

	    Within  Ottemo  you  can  mostly  meet  2   types  of  packages:  interfaces  and  implementations  (models  and
	actors).However it is not means that package can't play both roles - it rather exception than a rule.

	    Within  "interface/model" package  you can  meet "interfaces.go",  "manager.go", "helpers.go"  files. So,  these
	packages declares interfaces for future implementation packages,  but have no implementation inside. Them describing
	how  model should  look  like, but  have no  references  to any  model  candidate.  These  packages mostly  contains
	lightweight code.

	    "Implementation/actor" package  type represented by "init.go",  "decl.go", "helpers.go", "[unit].[i_interface]",
	"[unit].[other]" and other files. These packages comparatively  big, have external dependencies and most likely (but
	not necessarily), implements one or couple interfaces. Actor  packages using "init.go" to make self registration and
	should be closed to others packages. So, o

	    The necessity of having mentioned packages types follows  from Ottemo requirement to be pluggable and extendable
	as well  as GO language requirement  to not have  package cyclic dependencies. Last  one means that package  "A" can
	depend  on   package  "B",   but  package   "B"  can   not  depend  on   package  "A"   then  (for   any  subpackage
	respectively). Interfaces helps to  solve that issue as them contains interface declaration.  So, in previous sample
	if "A"  and "B" packages have  declarations, we can move  them to package "C"  and refer both to  that package, that
	package will have routines to  receive and send back objects which satisfies interface and  packages "A" and "B" can
	work to each other through  this interface functions and not knowing about each other directly.  As you can see this
	also allows to have replaceable packages - candidate on role.



3. Ottemo type naming convection:

	    GO language have a limited set of built  in types (ref https://golang.org/ref/spec#Types"), but have a method to
	declare any  amount of user  types. These  types can be  alias to base  type, struct,  function, etc. As  the naming
	convention is  same far all of  them, in Ottemo  we want to be  clear on object  types we are working  with, without
	making investigation on name.

	    Ottemo have  the following type name  agreements, perhaps them are  little bit extending type  names and violate
	naming recommendations, but it's worth it.

		  naming                 |  meaning
		-------------------------|--------------------------------------------------------
		Struct[TypeName]         | type [Name] struct { ... }
		Interface[InterfaceName] | type [Name] interface { ... }
		Func[FunctionName]       | type [Name] func { ... }
		[ObjectName]             | type [Name] struct { ... }; func (it *[Name]) { ... }


	    If you have an  GO language structure object, it can be  just structure to hold values together or  be a class -
	i.e.  to have methods. So,  there is slight difference between just structure and  structure with bound functions to
	it - methods  (ref https://golang.org/ref/spec#Method_declarations). First one (just structure)  have fields but you
	can't perform any action  on it - it just data holder,  and second one - have functions/methods you  can use to work
	with. First type will be named with "Struct" prefix in Ottemo and second one - without prefix.

	    Interface and function types are also very close. So  in first case you have a structure pointer with predefined
	set of functions  (if structure/class do not have  all of them it is  not interface capable - interface  not fits on
	it),  and in  second  case  you have  function  pointer  which also  will  not be  capable  if  will have  different
	descriptors. Both  of them are  pointers to  some objects or  nil value pointer,  and both  of them have  build time
	checking on interface appliance.

	    We want to be clear on all types of objects you will meet in code so prefixes will help you to know this without
	traveling a lot around surrounding files.


4. Variables / constants names:

	    GO language mostly have one fundamental principe of identifier naming - if the first character of the identifier
	name is a Unicode upper case  letter it will be exported, otherwise - not. Export -  means that you will have public
	visibility,  and  other   packages  will  be  able  to   use  that  item.   This  concerns  types   as  well.   (ref
	https://golang.org/ref/spec#Exported_identifiers)

	    In official documentation there  are no other restricts on this except first  letter which flags export. However
	in  other sources  you  can found  recommendation  to use  only  UpperCamelCase and  lowerCamelCase  naming for  all
	identifiers (ref https://code.google.com/p/go-wiki/wiki/CodeReviewComments).

	    Ottemo follows GO language official documentation in usage of first identifier letter and follows recommendation
	on camel-case naming of variable. However for a constants (to be good visible in code) "Const" prefix used.

*/
package main
