import axios from 'axios';
import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import { serializeAxiosError } from './reducer.utils';
import { apiUrl } from 'app/config/constants';

import { IUser } from 'app/shared/models/user.model';

const initialState = {
  loading: false,
  errorMessage: null,
  users: [] as ReadonlyArray<IUser>,
  account: {} as Readonly<IUser>,
};


// Async Actions
export const getUsers = createAsyncThunk('userManagement/fetch_users', async () => {
  const requestUrl = `${apiUrl}/users`;
  return axios.get<IUser[]>(requestUrl);
});

export const getAccount = createAsyncThunk('authentication/get_account', async () => {
  return axios.get<IUser>(`${apiUrl}/account`);
});


export type UserManagementState = Readonly<typeof initialState>;

export const UserManagementSlice = createSlice({
  name: 'userManagement',
  initialState: initialState as UserManagementState,
  reducers: {
    reset() {
      return initialState;
    },
  },
  extraReducers(builder) {
    builder
      .addCase(getUsers.pending, state => {
        state.loading = true;
      })
      .addCase(getUsers.rejected, (state, action) => ({
        ...state,
        loading: false,
        errorMessage: action.error.message,
      }))
      .addCase(getUsers.fulfilled, (state, action) => ({
        ...state,
        loading: false,
        users: action.payload.data,
      }))
      .addCase(getAccount.pending, state => {
        state.loading = true;
      })
      .addCase(getAccount.rejected, (state, action) => ({
        ...state,
        loading: false,
        errorMessage: action.error.message,
      }))
      .addCase(getAccount.fulfilled, (state, action) => ({
        ...state,
        loading: false,
        account: action.payload.data,
      }));
  },
});

export const { reset } = UserManagementSlice.actions;

// Reducer
export default UserManagementSlice.reducer;
