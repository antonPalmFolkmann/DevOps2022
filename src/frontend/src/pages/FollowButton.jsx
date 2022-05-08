import React, {useState, useContext} from 'react';
import Button from '@mui/material/Button';
import userContext from '../utils/userContext';
import { follow } from '../api/follow'
import { Container, Typography } from '@mui/material';

export default function FollowButton() {
    const user = useContext(userContext)

    return (
        <Container>
            <Button onClick={() => { follow(user.currentProfile.username).then((response) => {
                console.log(JSON.stringify(response))
            })}}>
                Follow
            </Button>
        </Container>    
    )
}