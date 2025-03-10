export class Channel {
  id
  rf_user_id
  channel_name
  live
  rf_active_stream_id
  channel_token
  icon
  created_at
  updated_at
  constructor(
    id: string,
    rf_user_id: string,
    channel_name: string,
    live: boolean,
    rf_active_stream_id: string,
    channel_token: string,
    icon: string,
    created_at: Date,
    updated_at: Date
  ) {
    this.id = id
    this.rf_user_id = rf_user_id
    this.channel_name = channel_name
    this.live = live
    this.rf_active_stream_id = rf_active_stream_id
    this.channel_token = channel_token
    this.icon = icon
    this.created_at = created_at
    this.updated_at = updated_at
  }
}

export class ChannelData {
  id
  channel_name
  icon
  constructor(id: string, channel_name: string, icon: string) {
    this.id = id
    this.channel_name = channel_name
    this.icon = icon
  }
}

export class Recording {
  id
  channel_id
  date
  duration
  channel_name
  path
  poster
  constructor(
    id: string,
    channel_id: string,
    date: Date,
    duration: number,
    channel_name: string,
    path: string,
    poster: string
  ) {
    this.id = id
    this.channel_id = channel_id
    this.date = date
    this.duration = duration
    this.channel_name = channel_name
    this.path = path
    this.poster = poster
  }
}

export class ChatMessage {
  stream_id
  text
  time
  channel
  constructor(stream_id: string, text: string, time: string, channel: ChannelData) {
    this.stream_id = stream_id
    this.text = text
    this.time = time
    this.channel = channel
  }
}
