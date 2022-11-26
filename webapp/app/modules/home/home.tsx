import { useAppDispatch, useAppSelector } from 'app/config/store';
import { IAlbum } from 'app/shared/models/album.model';
import React, {useEffect} from 'react';
import {
  EuiPage,
  EuiPageSection,
  EuiPageSidebar,
  EuiPageBody,
} from '@elastic/eui';

export const Home = () => {
  return (
    <EuiPage>
      <EuiPageSidebar paddingSize="l">
      </EuiPageSidebar>
      <EuiPageBody paddingSize="none" panelled="true">
        <EuiPageSection></EuiPageSection>
      </EuiPageBody>
    </EuiPage>
  );
}

export default Home;
