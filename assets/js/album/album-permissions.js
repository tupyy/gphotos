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
                <div class="col permission-user">` +
                    "<div class=\"fw-hold\">" + username + "</div>" +
                `</div>
                <div class="col permission-user">` +
                badges +
                `</div>
                <div class="col permission-remove-btn">
                    <button class="btn btn-danger btn-sm" onclick="removePermission(` + username +`)">Remove</button>
                </div>
            </div>
        </li>`
        );
    };
});
