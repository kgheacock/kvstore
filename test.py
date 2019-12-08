{
    "successful-responses": {





        "get-key": {
            "message"       : "Retrieved successfully",
            "doesExist"     : true,
            "value"         : "sampleValue",
            "address"       : "10.10.0.4:13800",
            "causal-context": {}
        },

        "delete-key": {
            "message"       : "Deleted successfully",
            "doesExist"     : true,
            "address"       : "10.10.0.3:13800",
            "causal-context": {}
        }

    },
    "failure-responses": {
        "insert-key-missing": {
            "message"       : "Error in PUT",
            "error"         : "Value is missing",
            "causal-context": {}
        },

        "insert-key-long": {
            "message"       : "Error in PUT",
            "error"         : "Key is too long",
            "address"       : "10.10.0.4:13800",
            "causal-context": {}
        },

        "update-key-missing": {
            "message"       : "Error in PUT",
            "error"         : "Value is missing",
            "causal-context": {}
        },

        "get-key": {
            "message"       : "Error in GET",
            "error"         : "Key does not exist",
            "doesExist"     : false,
            "address"       : "10.10.0.4:13800",
            "causal-context": {}
        },

        "delete-key": {
            "message"       : "Error in DELETE",
            "error"         : "Key does not exist",
            "doesExist"     : false,
            "causal-context": {}
        }
    }
}