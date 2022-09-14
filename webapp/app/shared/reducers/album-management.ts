import axios from 'axios';
import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import { apiUrl } from "app/config/constants";
import { IAlbum } from "../models/album.model";

const DEFAULT_PAGE_SIZE = 20;

const initialState = {
  loading: false,
  errorMessage: null,
  albums: [] as ReadonlyArray<IAlbum>,
  count: 0,
  offset: 0,
  limit: DEFAULT_PAGE_SIZE,
};

interface internalIAlbum {
  albums: IAlbum[];
  count: number;
}

export const getAlbums = createAsyncThunk('albumManagement/fetch_albums', async (offset, limit) => {
  const requestUrl = `${apiUrl}/albums?offset=${offset}&limit=${limit}`;
  return axios.get<internalIAlbum>(requestUrl);
});

export type AlbumManagementState = Readonly<typeof initialState>;

export const AlbumManagementSlice = createSlice({
  name: 'albumManagement',
  initialState: initialState as AlbumManagementState,
  reducers: {
    reset() {
      return initialState;
    },
  },
  extraReducers(builder) {
    builder
    .addCase(getAlbums.pending, state => {
      state.loading = true;
    })
    .addCase(getAlbums.rejected, (state, action) => ({
      ...state,
      loading: false,
      errorMessage: action.error.message,
    }))
    .addCase(getAlbums.fulfilled, (state, action) => {
      const d = action.payload.data;
      return {
        ...state,
        loading: false,
        albums: d.albums,
        count: d.count,
      }
    });
  },
});

export const { reset } = AlbumManagementSlice.actions;

export default AlbumManagementSlice.reducer;
