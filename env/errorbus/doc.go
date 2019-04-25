// Copyright 2019 Ottemo. All rights reserved.

/*
Package errorbus is a default implementation of InterfaceErrorBus declared in "github.com/ottemo/commerce/env" package.

Error bus is a service which provides a convenient way to deliver error to system administrator. Error bus uses logger
for an error message storing. Error message withing error bus provides with extra information  which allows to fileter
error messages. Following extra information used: "module name", "error level", "error code".

"module name" is a short-name for a package to have ability distinguish one package from another (no rules on it).

"error level" is one of "env" package constant. The more level of constant then it is less system related but more
application code related. So, system administrators most likely needs to monitor error messages with error level below 5.
Where as for customer applications should provide to user error messages above 5. Special meaning have error level 0 which
means that error happened outside Ottemo system package.

	env.ConstErrorLevelAPI        = 10
	env.ConstErrorLevelModel      = 9
	env.ConstErrorLevelActor      = 8
	env.ConstErrorLevelHelper     = 7
	env.ConstErrorLevelService    = 4
	env.ConstErrorLevelServiceAct = 3
	env.ConstErrorLevelCritical   = 2
	env.ConstErrorLevelStartStop  = 1
	env.ConstErrorLevelExternal   = 0

"error code" - unique string name identifying error  message. It is highly recommended to use "uuidgen" or similar tool
to get unique UUID identifier.

Once handled message will nto be handled second time it come to error bus as it flags as handled. In this way error
message can safely travel between application routines, knowing that it will be logged only once (first occurrence).

As "module name" and "error level" in most cases not changing for a particular package, it is recommended to use package
constants for them, like:
	ConstErrorModule = "rts"
	ConstErrorLevel  = env.ConstErrorLevelActor

There are two main usage approaches:

    Example 1: (Handling external errors)
    -------------------------------------
		if err := api.ValidateAdminRights(context); err != nil {
			return env.ErrorDispatch(err)
		}

    Example 2: (Generating error)
    -----------------------------
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "d4962872-9ee7-4c86-94ef-7325cb6e1c9d", "ID is not valid")
*/
package errorbus
