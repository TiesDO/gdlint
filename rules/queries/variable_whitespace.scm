[
    (variable_statement
        name: (_) @var.name
        "=" @var.assign
        value: (_) @var.value
    )
    (variable_statement
        name: (_) @var.name
        ":" @var.colon
        type: (_) @var.type
        (
            "=" @var.assign
            value: (_) @var.value
        )?
    )
    ; (variable_statement
    ;     name: (_) @var.name
    ;     ":=" @var.inferred
    ;     value: (_) @var.value
    ; )
]
