-- Create materialized view used to get all permissions for one album
CREATE MATERIALIZED VIEW album_user_permissions_view AS
    SELECT album.id, album.owner_id, album_user_permissions.user_id, album_user_permissions.permissions 
    FROM album
    JOIN album_user_permissions ON album.id = album_user_permissions.album_id;

-- Create materialized view used to albums' permissions by group
CREATE MATERIALIZED VIEW album_group_permissions_view AS
    SELECT album.id, album.owner_id, album_group_permissions.group_name, album_group_permissions.permissions
    FROM album
    JOIN album_group_permissions ON album.id = album_group_permissions.album_id;

-- Create trigger to update album_user_permissions_view
CREATE FUNCTION update_permissions_views() 
    RETURNS TRIGGER 
    LANGUAGE plpgsql
    SECURITY DEFINER
AS $$
BEGIN
    REFRESH MATERIALIZED VIEW album_user_permissions_view;
    REFRESH MATERIALIZED VIEW album_group_permissions_view;
    RETURN NULL;
END;
$$;

CREATE TRIGGER update_user_permissions_view_trigger AFTER UPDATE OR INSERT OR DELETE ON album_user_permissions FOR EACH ROW EXECUTE PROCEDURE update_permissions_views();
CREATE TRIGGER update_group_permissions_view_trigger AFTER UPDATE OR INSERT OR DELETE ON album_group_permissions FOR EACH ROW EXECUTE PROCEDURE update_permissions_views();

GRANT EXECUTE ON FUNCTION update_permissions_views TO core_readwrite;
GRANT SELECT ON album_user_permissions_view to core_readwrite;
GRANT SELECT ON album_group_permissions_view to core_readwrite;


