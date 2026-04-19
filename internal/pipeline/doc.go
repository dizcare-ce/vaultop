// Package pipeline provides a lightweight sequential execution engine for
// multi-step secret operations in vaultop.
//
// A Pipeline is composed of named Steps. Each Step receives a shared State
// that can be read and mutated by subsequent steps. Execution stops on the
// first error unless ContinueOnError is enabled.
//
// Built-in steps (FetchStep, ValidateStep, WriteStep) cover the common
// fetch → validate → write workflow, but custom StepFuncs can be added for
// any additional logic such as auditing, notifications, or caching.
package pipeline
