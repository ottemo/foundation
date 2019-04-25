// Copyright 2019 Ottemo. All rights reserved.

/*
Package utils provides a set of functions you should use to reduce and standardize repeatable cross packages code. So, it
considered that application code will use utils package as often as can.


"datatypes.go" - contains set of data-type declaration and conversion routines

  When you are writing an application code you should be aware of data-types application can work with.So, this file contains
  that information for you, in addition to set of types conversions routines. As utils package it is globally used byt all
  application parts it makes application to works a unified way on data-types. You should use this package routines to making
  types conversion to be sure other application parts will understand you.

  Example 1:
  ----------
        x := map[string]interface{} {"a": 10, "b": "20", "c": true")
        y := utils.InterfaceToString(x)
        z := utils.InterfaceToMap(y)
        fmt.Println(x, z)

        xTime := utils.InterfaceToTime("2010-01-15 10:22")
        xFloat := utils.InterfaceToFloat64("2", 10)
        xBool := utils.InterfaceToBool("yes")
        fmt.Println(xTime, xFloat, xBool)

  Example 2:
  ----------
        typeName := "[]string"
        typeValue := "1,2,3,4"

        typeInfo := utils.DataTypeParse( typeName )
        if typeInfo.IsKnown && typeInfo.IsArray  {
            typedValue :=  utils.InterfaceToArray( utils.StringToType(typeValue, typeName) )

            for _, value := typedValue {
                fmt.Println(value)
            }
        }


"generic.go" - contains set of unclassified routines should be unified due to application.

  Example:
  --------
      x := map[string]interface{} {"a": "10", "b": "20")

      if utils.KeysInMapAndNotBlank(x, "a", "b") {
          fmt.Println( utils.InterfaceToInt(x[a]) + Utils.InterfaceToInt(x[b]) )
          fmt.Println( utils.GetFirstMapValue(x) )
      }

      y := utils.RoundPrice( utils.InterfaceToFloat("11.22241") )
      fmt.Println(y)

      z := utils.SplitQuotedStringBy(`"a", "b"; "c,d", "e.f".`, ";,. ")
      fmt.Println(z)

      searchString := "just {sample} [code]"
      escapedSearchString := EscapeRegexSpecials(searchString)
      matched, err = regexp.MatchString(searchString, escapedSearchString)
      fmt.Println(matched, err)


"crypt.go" - provides an centralized way for bi-directional crypt of secure data.

  Notes:
      - SetKey() makes change for entire application. So, if you want local effect you should restore it after usage
      - normally application should take care about SetKey() on init and you should not touch it
      - if SetKey() was not called during application init then default hard-coded key will be used

  Example 1:
  ----------
      source := "just test"
      encoded := utils.EncryptStringBase64(source)
      decoded := utils.DecryptStringBase64(encoded)
      println( "'" + source + "' --encode--> '" + encoded + "' --decode--> '" +  decoded + "'")

      Output:
        'just test' --encode--> 'Ddryse1yNL5z' --decode--> 'just test'

  Example 2:
  ----------
      sampleData := []byte("It is just a sample.")

      outFile, _ := os.OpenFile("sample.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
      defer outFile.Close()
      writer, _ := utils.EncryptWriter(outFile)
      writer.Write(sampleData)

      inFile, _ := os.OpenFile("sample.txt", os.O_RDONLY, 0600)
      defer inFile.Close()
      reader, _ := utils.EncryptReader(inFile)
      readBuffer := make([]byte, 10)

      reader.Read(readBuffer)
      println(string(readBuffer))
      reader.Read(readBuffer)
      println(string(readBuffer))

      Output:
        It is just
         a sample.

"json.go" - contains set of json conversion related routines

  Example:
  --------
      x := map[string]interface{} {"a": "10", "b": "20"}
      y := utils.EncodeToJSONString(x)
      fmt.Println(y)


"templates.go" - endpoint to register new application scope template functions as well simplifies GO templates usage.
In following sample "first" directive will be available to any code addressing text templates to utils.TextTemplate(...)

  Example:
  --------
      newFunc := func(identifier string, args ...interface{}) string {
          if len(args) > 0 {
              return utils.InterfaceToString(args[0])
          }
          return "null"
      }
      utils.RegisterTemplateFunction("first", newFunc)

      context := map[string]interface{} {"a": 10, "label": "first argument: "}
      result, err := utils.TextTemplate(".label {{first .a}}", context)
*/
package utils
