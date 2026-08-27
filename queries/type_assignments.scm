; untyped variable declarations
(variable_statement
    !type
    name: (name) @untyped_variable_statement
)

; untyped function return values
(function_definition
    !return_type
    name: (name) @untyped_function_return
)

; untyped function arguments
(function_definition
    parameters: (parameters
        (identifier) @untyped_function_argument
    )
)
