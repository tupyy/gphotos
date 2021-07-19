$(document).ready(function() {
    let permissions = {};

    getPermission = function() {
        let result = {};
        result['permissions'] = [];

        result['username'] = $("#select-users option:selected").text();
        result['index'] = $("#select-users").val();

        $(".form-check input").each(function () {
            if ( $(this).is(":checked") ) {
                result.permissions.push($(this).val());
            }
        });

        if (permissions.hasOwnProperty(result.username)) {
            throw 'permission already added for this user';
        }

        permissions[result.username] = result
        return result;
    }

    $("#userPermissionButton").click(function () {
        $(".row .no-permission").remove();

        try {
            p = getPermission();
            addPermissionElement(p);
        } catch (e) {
            console.log(e);
        }
    });

    let addPermissionElement = function(permission) {
        let badges = "";

        permission.permissions.forEach(function(item) {
            badges += '<span class="badge bg-success">' + item + "</span>" 
        });

        $("#selected-users").append(`
        <li class="list-group-item">
            <div class="row">
                <div class="col permission-user">` +
                    "<div class=\"fw-hold\">" + permission.username + "</div>" +
                `</div>
                <div class="col permission-user">` +
                badges +
                `</div>
                <div class="col permission-remove-btn">
                    <button class="btn btn-danger btn-sm">Remove</button>
                </div>
            </div>
        </li>`
        );
    };
});
