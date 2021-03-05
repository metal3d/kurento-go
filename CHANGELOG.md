# Changelog

-  Added GetPipelines functionality in relation to ServerManager
-  Remove `connections` map as it was unused and causing concurrent map write panics
-  Fix goroutine management to match conventions
-  Return error from connection.Request func to surface websocket errors correctly