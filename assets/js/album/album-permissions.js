$(() => {
    let uPermissions = {};
    let gPermissions = {};

    $("#userPermissionButton").on("click", function () {
        try {
            p = getPermission(".container-permissions-user", "#select-users");
            if (!uPermissions.hasOwnProperty(p.id)) {
                if (p.permissions.length > 0) { 
                    $(".row .no-users-permission").remove();
                    uPermissions[p.id] = p.permissions;
                    addPermissionElement("#selected-users",p.username,p.id, p.permissions);
                }       
            } 
            $("#inputUserPermissions").attr('value',JSON.stringify(uPermissions));
        } catch (e) {
            console.log(e);
        }
    });
    
    $("#groupPermissionButton").on("click", function () {
        try {
            p = getPermission(".container-permissions-group", "#select-groups");
            if (!gPermissions.hasOwnProperty(p.username)) {
                if (p.permissions.length > 0) {
                    $(".row .no-groups-permission").remove();
                    gPermissions[p.username] = p.permissions;
                    addPermissionElement("#selected-groups",p.username, p.username, p.permissions);
                }
            }
            $("#inputGroupPermissions").attr('value', JSON.stringify(gPermissions));
        } catch (e) {
            console.log(e);
        }
    });
    
    $("#selected-users").on("click",'.remove-permission', function (e) {
        e.preventDefault();

        let parent = $(this).parents("li")
        let username = $(parent).find("input").val();

        if (uPermissions.hasOwnProperty(username)) {
            delete uPermissions[username];
            $("#inputUserPermissions").attr('value',JSON.stringify(uPermissions));
        }

        if (Object.keys(uPermissions).length === 0) {
           $(".container-permissions-selected-users").append(`
                <div class="title no-users-permission">None</div>
            `); 
        }

        $(parent).remove();

    });
    
    $("#selected-groups").on("click",'.remove-permission', function (e) {
        e.preventDefault();

        let parent = $(this).parents("li")
        let name = $(parent).find("input").val();

        if (gPermissions.hasOwnProperty(name)) {
            delete gPermissions[name];
            $("#inputGroupPermissions").attr('value', JSON.stringify(gPermissions));
        }

        if (Object.keys(gPermissions).length === 0) {
           $(".container-permissions-selected-groups").append(`
                <div class="title no-groups-permission">None</div>
            `); 
        }

        $(parent).remove();

    });

    const getPermission = function(parent, element) {
        let permissions = [];

        id = $(element + " option:selected").val();
        username = $(element + " option:selected").text();
        

        $(parent + " .form-check input").each(function () {
            if ( $(this).is(":checked") ) {
                permissions.push($(this).val());
            }
        });

        return {id:id, username: username, permissions: permissions};
    }

    const addPermissionElement = function(dest, username, id, permissions) {
        let badges = "";

        permissions.forEach(function(item) {
            badges += '<span class="badge bg-success">' + item + "</span>" 
        });

        $(dest).append(`
        <li class="list-group-item">
            <div class="row">
                <input type="hidden" value="` + id + `"/>
                <div class="col permission-user">` +
                    "<div>" + username + "</div>" +
                `</div>
                <div class="col permission-user">` +
                badges +
                `</div>
                <div class="col permission-remove-btn">
                    <button class="btn btn-outline-danger btn-sm remove-permission">Remove</button>
                </div>
            </div>
        </li>`
        );
    };


    if ( $("#inputUserPermissions").val() !== "" ) {
        uPermissions = JSON.parse($("#inputUserPermissions").val());
        
        Object.keys(uPermissions).forEach((k) => {
            let username = "";
            $("#select-users > option").each(function() {
                if ( this.value === k ) {
                    username = this.text;
                }
            });
            // add div element
            addPermissionElement("#selected-users", username, k, uPermissions[k]);
        }); 
    }

    if ( $("#inputGroupPermissions").val() !== "" ) {
        try {
            gPermissions = JSON.parse($("#inputGroupPermissions").val())
        
            Object.keys(gPermissions).forEach((k) => {
                // add div element
                addPermissionElement("#selected-groups", k, k, gPermissions[k]);
            })
        } catch(err) {
            console.log(err);
        }
    }
});
