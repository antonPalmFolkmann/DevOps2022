import React, { useContext, useEffect} from 'react';
import Box from '@mui/material/Box';
import List from '@mui/material/List';
import Card from '@mui/material/Card';
import Stack from '@mui/material/Stack';
import CardContent from '@mui/material/CardContent';
import Typography from '@mui/material/Typography';
import { Container } from '@mui/material';
import Profile from './Profile';
import Button from '@mui/material/Button';
import NavBar from './NavBar';
import userContext from '../utils/userContext';
import { getMessages, exportMessages, exportProfile, getMessagesByUser } from '../api/messages'
import { logout } from '../api/logout' 
import { getTimestamp } from '../utils/functions';

export default function Timeline() {
  const user = useContext(userContext)
  useEffect(() => {
    getMessages()
    .then((response) => {
      user.setCurrentMessages(exportMessages)
    })
  }, [])
  useEffect(() => {
    logout()
    .then((response) => {})
  }, [])
  let screenHeight = window.innerHeight;

  return (
    <Container 
    sx={{
      display: 'flex',
    }}>
    {(() => {
    if (!(user.username !== '' && user.currentProfile.username === '')) {
    return (
      <Profile />
    )
    }
    })()}

      <Container style={{
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        paddingTop: 96,
        maxHeight: screenHeight,
        minWidth: 630,
        overflowY: 'auto',
        overflowX: 'hidden', 
      }}>
        <List style={{
          maxHeight: screenHeight
        }}>
            {user.currentMessages.map( (msg) => {
              return (
                <Box>
                  <Card variant="outlined" sx={{ minWidth: 600, borderColor: '#a3c9fe'}} onClick={() => { 
                      if (user.username !== '') {
                        getMessagesByUser(msg.authorName)
                        .then((response) => {
                          console.log('Parsed messages in Timeline');
                          user.setCurrentProfile(exportProfile);
                          user.setCurrentMessages(exportMessages);
                        })
                      } 
                    }}>
                    <CardContent>
                      <Stack direction="column">                  
                          <Typography variant="h5" component="div">
                            {msg.text} 
                          </Typography>
                          <Stack direction="row">
                          <Typography variant="p" component="div">
                            {msg.authorName + ' /'} 
                          </Typography>
                          <Typography variant="p" component="div">
                            {'/ ' + getTimestamp(msg.pubDate)} 
                          </Typography>
                          </Stack>
                      </Stack>
                    </CardContent>
                  </Card>
                  <p>
                    {'\n'}
                  </p>
                </Box>
              )
            })}
        </List>
    </Container>
    {(() => {
    if (user.username !== '') {
    return (
      <NavBar/>
    )
    }
    })()}
  </Container>      
  );
}
