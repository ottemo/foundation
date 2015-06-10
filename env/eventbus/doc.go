// Copyright 2015 The Ottemo Authors. All rights reserved.

/*
Package eventbus is a default implementation of InterfaceEventBus declared in "github.com/ottemo/foundation/env" package.

Event bus is a service used for simplified communication between application code. Event provider emits an event message
and event listeners makes special handling for an event.

Event name is "." delimited string. So, even listeners can listen for all messages of "top level" message (i.e. listener
for "api" event will listen for "api.checkout.visitCheckout" automatically).

Event provides a data objects relative to. These objects could be changed during event handling, as well as new data
could be added to a data map, during event processing.

To be more consistent and clear, event names should be declared as a package constants with description about providing
event data map.

    Example 1:
    ----------
        return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "004e9f7b-bb97-4356-bbc2-5e084736983b", "unknown cmd '"+args[0]+"'")
        env.Event("api.checkout.visitCheckout", eventData)

    Example 2:
    ----------
        salesHandler := func(event string, eventData map[string]interface{}) bool {
            env.LogMessage( fmt.Sprintf("%+v", eventData) )
        }
        env.EventRegisterListener("checkout.success", salesHandler)
*/
package eventbus
