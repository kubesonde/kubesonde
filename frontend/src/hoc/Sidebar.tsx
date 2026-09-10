import React, { useState } from "react";
import {
  ProSidebar,
  Menu,
  MenuItem,
  SidebarHeader,
  SidebarContent,
  SidebarFooter,
} from "react-pro-sidebar";
import { useLocation, useNavigate } from "react-router-dom";
import { FiHome, FiArrowLeftCircle, FiArrowRightCircle } from "react-icons/fi";
import { FaGithub, FaList } from "react-icons/fa";

//menuCollapse state using useState hook
import { GrGraphQl } from "react-icons/gr";

const routes = [
  /*{
         name: 'Fetch from Kubesonde instance',
         path: '/remote'
     },*/
  {
    name: "Upload a file",
    path: "/upload",
    icon: <FaList />,
  },
  {
    name: "Load example probe",
    path: "/example",
    icon: <GrGraphQl />,
  },
];

export const Sidebar: React.FC = () => {
  const [menuCollapse, setMenuCollapse] = useState<boolean>(false);
  const { pathname } = useLocation();
  const navigate = useNavigate();
  const menuIconClick = () => {
    menuCollapse ? setMenuCollapse(false) : setMenuCollapse(true);
  };
  const MenuItems = routes.map((route) => (
    <MenuItem
      key={route.name}
      icon={route.icon}
      className="menuItem"
      active={pathname === route.path}
      onClick={() => navigate(route.path)}
    >
      {route.name}
    </MenuItem>
  ));
  return (
    <>
      <div id="header">
        {/* collapsed props to change menu size using menucollapse state */}
        <ProSidebar collapsed={menuCollapse} role="sidebar">
          <SidebarHeader>
            <div className="logotext">
              <img src="/logo257.png" alt="" className="brandmark" />
              <p>{menuCollapse ? "Ksonde" : "Kubesonde Viewer"}</p>
            </div>
            <div className="closemenu" onClick={menuIconClick}>
              {/* changing menu collapse icon on click */}
              {menuCollapse ? <FiArrowRightCircle /> : <FiArrowLeftCircle />}
            </div>
          </SidebarHeader>
          <SidebarContent>
            <Menu iconShape="square">
              <MenuItem
                icon={<FiHome />}
                active={pathname === "/"}
                onClick={() => navigate("/")}
              >
                Home
              </MenuItem>
              {MenuItems}
            </Menu>
          </SidebarContent>
          <SidebarFooter>
            <Menu iconShape="square">
              <MenuItem className="version">
                Version {import.meta.env.VITE_APP_VERSION}
              </MenuItem>
              <MenuItem
                icon={<FaGithub />}
                onClick={() =>
                  window.open(
                    "https://github.com/kubesonde/kubesonde",
                    "_blank",
                    "noopener,noreferrer"
                  )
                }
              >
                View on GitHub
              </MenuItem>
            </Menu>
          </SidebarFooter>
        </ProSidebar>
      </div>
    </>
  );
};
