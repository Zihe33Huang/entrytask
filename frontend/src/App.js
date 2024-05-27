import './App.css';
import LoginPage from "./pages/Login";
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import ProfilePage from "./pages/Profile";
import {NotificationContainer} from "react-notifications";
import React from "react";


function App() {
    return (
        <Router>
            <Routes>
                <Route path="/" element={<LoginPage />} />
                <Route path="/profile" element={<ProfilePage />} />
            </Routes>
            <NotificationContainer/>
        </Router>

    );
}

export default App;
