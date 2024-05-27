import React, {useEffect, useState} from 'react';
import { FaEdit } from 'react-icons/fa';
import "./InpurTextField.css"
import {FaArrowPointer} from "react-icons/fa6";

const EditableTextField = ({ initialValue }) => {
    const [nickname, setNickname] = useState(initialValue);
    const [isEditing, setIsEditing] = useState(false);

    // useEffect(() => {
    //     console.log('hhhh' + value)
    // }, [value]);
    const handleEditClick = () => {
        setIsEditing(true);
        console.log(nickname)
    };


    const handleChange = (e) => {
        setNickname(e.target.value);
    };

    const handleSave = () => {
        setIsEditing(false);
        // Here you can add code to save the edited value to your database or perform any other necessary actions
    };

    return (
        <div className="nickname">
            {isEditing ? (
                <input
                    type="text"
                    value={nickname}
                    onChange={handleChange}
                    onBlur={handleSave}
                    autoFocus
                />
            ) : (
                <p>Nickname：<span>{nickname}</span></p>
            )}
            {!isEditing && (
                <button onClick={handleEditClick}>
                    <FaEdit />
                </button>
            )}
            {isEditing && (
                <button onClick={handleEditClick}>
                    <FaArrowPointer />
                </button>
            )}
        </div>
    );
};

export default EditableTextField;