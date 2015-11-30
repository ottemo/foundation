// Copyright 2016 Ottemo, All rights reserved.

/*
Package cron is a utility to schedule tasks.  These tasks maybe scheduled for
a specific time, they may be repeatable or intended to be run immediately.
There are several API Endpoints which allow a developer to interact with the
scheduling system.

It is important to understand several concepts.

<b>Task</b> - a job which can be scheduled to run at a specific time
<b>Schedule</b> - a listing of all active tasks, when they will be executed and their respective metadata

The API allows you to:
        * Obtain a list of the currently scheduled tasks
        * Create a task to be run on a schedule
        * Obtain a list of possible tasks to be scheduled
        * Enable a task to be run on a schedule
        * Disable a task
        * Update the specified task
        * Run the specified task now

//TODO: add link to api documentation

*/
package cron
