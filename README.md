ipupdater
---------

Simple golang program to update a route53 'A' record with the current IP of
the calling host. Written to handle the case where local IP tends to change,
either due to DHCP or instance recreations.

