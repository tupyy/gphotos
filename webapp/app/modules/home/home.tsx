import { useAppDispatch, useAppSelector } from 'app/config/store';
import { IAlbum } from 'app/shared/models/album.model';
import { getAlbums } from 'app/shared/reducers/album-management';
import React, {useEffect} from 'react';
import Album from './album';
import { Grid } from '@mui/material';

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
      <Grid container rowSpacing={0.5} columnSpacing={0.5}>
        {albumState.albums && !albumState.loading
          ? albumState.albums.map((album: IAlbum, index) => (
            <Grid item>
              <Album key={index} album={album}/>
            </Grid>
          ))
          : null
        }
      </Grid>
    </div>
  );
}

export default Home;
