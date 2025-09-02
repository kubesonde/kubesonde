import React, { useState } from "react";
import {
  ProSidebar,
  Menu,
  MenuItem,
  SidebarHeader,
  SidebarContent,
  SidebarFooter,
} from "react-pro-sidebar";
import { Link } from "react-router-dom";
import { FiHome, FiArrowLeftCircle, FiArrowRightCircle } from "react-icons/fi";

//menuCollapse state using useState hook
import { GrGraphQl } from "react-icons/gr";

const routes = [
  /*{
        name: 'Load from file',
        path: '/fileUpload',
        icon: <FaList/>
    },
     {
         name: 'Fetch from Kubesonde instance',
         path: '/remote'
     },*/
  {
    name: "Load example probe",
    path: "/example",
    icon: <GrGraphQl />,
  },
];

export const Sidebar: React.FC = () => {
  const [menuCollapse, setMenuCollapse] = useState<boolean>(false);
  const menuIconClick = () => {
    menuCollapse ? setMenuCollapse(false) : setMenuCollapse(true);
  };
  const MenuItems = routes.map((route) => (
    <MenuItem key={route.name} icon={route.icon} className="menuItem">
      <Link
        id={route.name}
        key={route.name}
        role="menuLink"
        className="menuOption"
        to={route.path}
      >
        {route.name}
      </Link>
    </MenuItem>
  ));
  return (
    <>
      <div id="header">
        {/* collapsed props to change menu size using menucollapse state */}
        <ProSidebar collapsed={menuCollapse} role="sidebar">
          <SidebarHeader>
            <div className="logotext">
              {/* Icon change using menucollapse state */}
              <p>{menuCollapse ? "Ksonde" : "Kubesonde Viewer"}</p>
            </div>
            <div className="closemenu" onClick={menuIconClick}>
              {/* changing menu collapse icon on click */}
              {menuCollapse ? <FiArrowRightCircle /> : <FiArrowLeftCircle />}
            </div>
          </SidebarHeader>
          <SidebarContent>
            <Menu iconShape="square">
              <MenuItem icon={<FiHome />}>
                <Link to="/">Home</Link>
              </MenuItem>
              {MenuItems}
            </Menu>
          </SidebarContent>
          <SidebarFooter>
            <Menu iconShape="square">
              <MenuItem>Version: {import.meta.env.REACT_APP_VERSION}</MenuItem>
            </Menu>
          </SidebarFooter>
        </ProSidebar>
      </div>
    </>
  );
};
