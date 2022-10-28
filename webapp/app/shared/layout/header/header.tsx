import { useAppDispatch, useAppSelector } from 'app/config/store';
import { getAccount } from 'app/shared/reducers/user-management';
import React, { useState } from 'react';
import {
  EuiAvatar,
  EuiFlexGroup,
  EuiFlexItem,
  EuiHeader,
  EuiHeaderLogo,
  EuiHeaderSectionItem,
  EuiHeaderSectionItemButton,
  EuiHeaderSection,
  EuiLink,
  EuiHeaderLink,
  EuiHeaderLinks,
  EuiPopover,
  EuiSpacer,
  EuiText,
  EuiSearchBar,
  useGeneratedHtmlId,
} from '@elastic/eui';
import { IUser } from 'app/shared/models/user.model';

const Header = () => {
  let dummyUser :IUser = {
    name: "test test",
    username: "test",
    role: "admin"
  }
  
  return (
  <EuiHeader>
    <EuiHeaderSection grow={false}>    
      <EuiHeaderSectionItem border="right">
        <EuiHeaderLogo>Photos</EuiHeaderLogo>
      </EuiHeaderSectionItem>
    </EuiHeaderSection>

    <EuiHeaderSection side="right" grow={false}>
      <EuiHeaderSectionItem>
        <EuiHeaderLinks aria-label="App navigation links example">
          <EuiHeaderLink isActive>Album nou</EuiHeaderLink>
        </EuiHeaderLinks>
      </EuiHeaderSectionItem>
      
      <EuiHeaderSectionItem>
        <HeaderUserMenu user={dummyUser}/>
      </EuiHeaderSectionItem>
    </EuiHeaderSection>
    </EuiHeader>
  )
}

const SearchBar = () => {
  return (
      <EuiSearchBar
        box={{
          placeholder: 'name:home -is:active joe',
        }}
      />
    );
}
interface IHeader {
  user: IUser,
}

const HeaderUserMenu = (props: IHeader) => {
  const headerUserPopoverId = useGeneratedHtmlId({
    prefix: 'headerUserPopover',
  });
  const [isOpen, setIsOpen] = useState(false);

  const onMenuButtonClick = () => {
    setIsOpen(!isOpen);
  };

  const closeMenu = () => {
    setIsOpen(false);
  };

  const button = (
    <EuiHeaderSectionItemButton
      aria-controls={headerUserPopoverId}
      aria-expanded={isOpen}
      aria-haspopup="true"
      aria-label="Account menu"
      onClick={onMenuButtonClick}
    >
      <EuiAvatar name={props.user.name} size="s" />
    </EuiHeaderSectionItemButton>
  );

  return (
    <EuiPopover
      id={headerUserPopoverId}
      button={button}
      isOpen={isOpen}
      anchorPosition="downRight"
      closePopover={closeMenu}
      panelPaddingSize="none"
    >
      <div style={{ width: 320 }}>
        <EuiFlexGroup
          gutterSize="m"
          className="euiHeaderProfile"
          responsive={false}
        >
          <EuiFlexItem grow={false}>
            <EuiAvatar name={props.user.name} size="xl" />
          </EuiFlexItem>

          <EuiFlexItem>
            <EuiText>
              <p>{props.user.name}</p>
            </EuiText>

            <EuiSpacer size="m" />

            <EuiFlexGroup>
              <EuiFlexItem>
                <EuiFlexGroup justifyContent="spaceBetween">
                  <EuiFlexItem grow={false}>
                    <EuiLink>Log out</EuiLink>
                  </EuiFlexItem>
                </EuiFlexGroup>
              </EuiFlexItem>
            </EuiFlexGroup>
          </EuiFlexItem>
        </EuiFlexGroup>
      </div>
    </EuiPopover>
  );
};
export default Header;
