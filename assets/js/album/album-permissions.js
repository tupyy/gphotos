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
                } else {
                    
                }       
            } 
        } catch (e) {
            console.log(e);
        }
    });
    
    $("#groupPermissionButton").on("click", function () {
        try {
            p = getPermission(".container-permissions-group", "#select-groups");
            if (!gPermissions.hasOwnProperty(username) && p.permissions.length > 0) {
                $(".row .no-groups-permission").remove();
                gPermissions[p.username] = p.permissions;
                addPermissionElement("#selected-groups",p.username, p.permissions);
            }
        } catch (e) {
            console.log(e);
        }
    });

    $("#create-album").submit(function(event) {
        event.preventDefault();
            jsonData= {
                name: $("#create-album #name").val(),
                description: $("#create-album #description").val(),
                location: $("#create-album #location").val(),
                user_permissions: uPermissions
            }

        $.ajax({
            url: "/album",
            type: "POST",
            dataType: "json",
            contentType: 'application/json; charset=utf-8',
            data: JSON.stringify(jsonData)
        }).done(function() {
            console.log("done");
        }).fail(function(e) {
            console.log(e);
        });
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


    const removePermission = (id) => {
        $("#"+id).remove();
    }

    const addPermissionElement = function(dest, username, permissions) {
        let badges = "";

        permissions.forEach(function(item) {
            badges += '<span class="badge bg-success">' + item + "</span>" 
        });

        $(dest).append(`
        <li class="list-group-item" id="` + username + `">
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
