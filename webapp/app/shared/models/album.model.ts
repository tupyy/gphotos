import { ITag } from "./tag.model";

type PermissionType = {
  [userID: string]: string[],
}

export interface IMedia {
  mediaType: string;
  filename: string;
  bucket: string;
  thumbnail: string;
  metadata: Map<string,string>;
  CreateDate: Date;
}

export interface IAlbum {
  id: string;
  name: string;
  date: Date;
  description: string;
  location: string;
  thumbnail: string;
  userPermissions: PermissionType;
  groupPermissions: PermissionType;
  Photos: IMedia;
  Videos: IMedia;
  tags: ITag[];
}
