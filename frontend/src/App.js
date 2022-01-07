import React from "react";
import {
  Routes,
  Route
} from "react-router-dom";

import Landing from "./pages/landing_page/Landing";
import Home from "./pages/home_page/Home";
import Auth from "./pages/auth_page/Auth";

const App = () => {
  return (
    <Routes>
      <Route exact path="/" element={<Landing />} />
      <Route path="/home" element={<Home />} />
      <Route path="/auth" element={<Auth />} />
    </Routes>
  );
}

export default App;
