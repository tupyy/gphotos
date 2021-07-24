$(() => {
    let uPermissions = {};
    let gPermissions = {};

    $("#userPermissionButton").on("click", function () {
        try {
            p = getPermission(".container-permissions-user", "#select-users");
            if (!uPermissions.hasOwnProperty(username)) {
                if (p.permissions.length > 0) { 
                    $(".row .no-users-permission").remove();
                    uPermissions[p.username] = p.permissions;
                    addPermissionElement("#selected-users",p.username, p.permissions);
                }       
            } 
            setPermissionInputValue("#inputUserPermissions", uPermissions);
        } catch (e) {
            console.log(e);
        }
    });
    
    $("#groupPermissionButton").on("click", function () {
        try {
            p = getPermission(".container-permissions-group", "#select-groups");
            if (!gPermissions.hasOwnProperty(username)) {
                if (p.permissions.length > 0) {
                    $(".row .no-groups-permission").remove();
                    gPermissions[p.username] = p.permissions;
                    addPermissionElement("#selected-groups",p.username, p.permissions);
                }
            }
            setPermissionInputValue("#inputGroupPermissions", gPermissions);
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
            setPermissionInputValue("#inputUserPermissions", uPermissions);
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
            setPermissionInputValue("#inputGroupPermissions", gPermissions);
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

        username = $(element + " option:selected").text();

        $(parent + " .form-check input").each(function () {
            if ( $(this).is(":checked") ) {
                permissions.push($(this).val());
            }
        });

        return {username: username, permissions: permissions};
    }

    const setPermissionInputValue = (inputID, permMap) => {
        let p = "";

        Object.keys(permMap).forEach(function(k) {
            let pp = "("+k+"#";

            permMap[k].forEach(function(item) {
                pp += item + ",";
            });

            pp = pp.slice(0,-1);
            p += pp + ")";
        });

        $(inputID).attr("value",p);
    }

    const addPermissionElement = function(dest, username, permissions) {
        let badges = "";

        permissions.forEach(function(item) {
            badges += '<span class="badge bg-success">' + item + "</span>" 
        });

        $(dest).append(`
        <li class="list-group-item">
            <div class="row">
                <input type="hidden" value="` + username + `"/>
                <div class="col permission-user">` +
                    "<div>" + username + "</div>" +
                `</div>
                <div class="col permission-user">` +
                badges +
                `</div>
                <div class="col permission-remove-btn">
                    <button class="btn btn-danger btn-sm remove-permission">Remove</button>
                </div>
            </div>
        </li>`
        );
    };
});
