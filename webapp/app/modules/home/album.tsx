import { IAlbum } from 'app/shared/models/album.model';
import React from 'react';
import Box from '@mui/material/Box';
import { Grid, Paper } from '@mui/material';
import { styled } from '@mui/material/styles';
import Card from '@mui/material/Card';
import CardHeader from '@mui/material/CardHeader';
import CardMedia from '@mui/material/CardMedia';
import CardContent from '@mui/material/CardContent';
import CardActions from '@mui/material/CardActions';
import Collapse from '@mui/material/Collapse';
import Avatar from '@mui/material/Avatar';
import IconButton, { IconButtonProps } from '@mui/material/IconButton';
import Typography from '@mui/material/Typography';
import { red } from '@mui/material/colors';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import dayjs from 'dayjs';
import { ALBUM_DATE_FORMAT } from 'app/config/constants';

interface IAlbumProps {
  album: IAlbum;
}

const Album = (props: IAlbumProps) => {
  return (
    <Card sx={{ maxWidth: 345 }}>
      <CardHeader
        avatar={
          <Avatar sx={{ bgcolor: red[500] }} aria-label="recipe">
            RC
          </Avatar>
        }
        action={
          <IconButton aria-label="settings">
            <MoreVertIcon />
          </IconButton>
        }
        title={ props.album.location }
        subheader={ dayjs(props.album.date).format(ALBUM_DATE_FORMAT) }
      />
      <CardMedia
        component="img"
        height="194"
        image={props.album.thumbnail}
        alt="album cover"
      />
      <CardContent>
        <Typography variant="body2" color="text.secondary">
          {props.album.name}
        </Typography>
      </CardContent>
    </Card>
  );
}

export default Album;
