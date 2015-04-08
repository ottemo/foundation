// Copyright 2014 The Ottemo Authors. All rights reserved.

/*
Package config is a default implementation of InterfaceConfig declared in "github.com/ottemo/foundation/env" package.

Default config values service is a database based storage of application configuration values. Is using one collection
within current DB storage service to do this. Config using [StructConfigItem] structure to represent its values. Except
key/value information it contains a lot more information to describe config value. This extra information supposed for
using within frontend editors, so it consist of:

    "description" - config value information
    "editor"      - coma separated editors frontend can use for value edit
    "label"       - readable config value name
    "options"     - any value(s) related to editor
    "type"        - config value type (refer to utils package)
    "image"       - some icon for config value (like paypal or visa icon, blank for most cases)

In addition to this information struct contains "type" value which using during DB to GO type conversion. Refer to utils
package on possible values application understands. It needed as some db service in most cases type based, so in SQLite
engine config value stored as text value and before usage it should be converted to appropriate type.

There is special type [env.ConstConfigTypeGroup] which represents top group for some config group, so it just a special
type application using to filter groups. For an instance path "app.checkout" could be a group then paths like
"app.checkout.allow_oversell" and "app.checkout.oversell_limit" are config values related to this group. So, with usage
of functions config.GetGroupItems() and config.GetItemsInfo() you can filter group based items.

Each config value can have validator function associated, which can also modify value puring validation.

To be more consistent and clear it is highly recommended to declare config value paths as a package constants.

    Example 1:
    ----------
    	const ConstConfigPathAllowOversell = "checkout.allow_oversell"

	    config := env.GetConfig()
	    if config != nil {
			validator := func(value interface{}) (interface{}, error) {
				return utils.InterfaceToBool(value), nil
			}

	        err := config.RegisterItem(env.StructConfigItem {
	            Path:        ConstConfigPathAllowOversell,
	            Value:       false,
	            Type:        env.ConstConfigTypeBoolean,
	            Editor:      "boolean",
	            Options:     nil,
	            Label:       "Allow oversell",
	            Description: "Allows oversell for out of stock items",
	            Image:       "",
	        }, validator)

	        if err != nil {
	        	return env.ErrorDispatch(err)
	        }
	    }

    Example 2:
    ----------
        if utils.InterfaceToBool( env.ConfigGetValue(checkout.ConstConfigPathAllowOversell) ) {
        	...
        }
*/
package config
