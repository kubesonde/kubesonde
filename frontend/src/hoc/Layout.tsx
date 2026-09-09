import {Sidebar} from "./Sidebar";
import {Route, Routes} from "react-router-dom";
import {GraphJSONUploadComponent} from "../components/graph/GraphJSONUploadComponent";
import './Layout.css'
import {HomeComponent} from "src/hoc/HomeComponent";
import {GraphFromLocation} from "src/components/graph/GraphFromLocation";
import {ExampleGraphComponent} from "src/components/graph/ExampleGraph";
import {GraphFromApi} from "src/components/graph/GraphFromApi";
import {isApiMode} from "src/utils/config";

export function Layout(){
    const MenuRoutes = (
        <Routes>
            <Route key="1" path="/example" element={<ExampleGraphComponent />}/>
            <Route key="2" path="/fileUpload" element={<GraphJSONUploadComponent />}/>
            <Route key="4" path="/upload" element={<HomeComponent/>}/>
            <Route key="0" path="/" element={isApiMode ? <GraphFromApi/> : <HomeComponent/>}/>
            <Route key="3" path="/graph/:graphName" element={<GraphFromLocation/>}/>
        </Routes>
    )
    return (
        <div className="app">
            <Sidebar />
            <div className="main">
                {MenuRoutes}
            </div>
        </div>
    )
}
