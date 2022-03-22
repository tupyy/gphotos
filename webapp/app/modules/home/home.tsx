import { useAppDispatch, useAppSelector } from 'app/config/store';
import { IAlbum } from 'app/shared/models/album.model';
import { getAlbums } from 'app/shared/reducers/album-management';
import React, {useEffect} from 'react';

export const Home = () => {
  const dispatch = useAppDispatch();
  const albumState = useAppSelector(state => state.albumManagement);

  useEffect(() => {
    dispatch(getAlbums());
  }, []);

  return (
    <div>
      {albumState.loading 
        ? (<div>Loading</div>)
        : null
      }
      <ul>
        {albumState.albums && !albumState.loading
          ? albumState.albums.map((album: IAlbum, index) => (
            <li key={index}>{album.name}</li>
          ))
          : null
        }
      </ul>
    </div>
  );
}

export default Home;
