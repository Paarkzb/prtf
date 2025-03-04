export class Channel {
  id
  rf_user_id
  channel_name
  live
  rf_active_stream_id
  channel_token
  created_at
  updated_at
  constructor(
    id: string,
    rf_user_id: string,
    channel_name: string,
    live: boolean,
    rf_active_stream_id: string,
    channel_token: string,
    created_at: Date,
    updated_at: Date
  ) {
    this.id = id
    this.rf_user_id = rf_user_id
    this.channel_name = channel_name
    this.live = live
    this.rf_active_stream_id = rf_active_stream_id
    this.channel_token = channel_token
    this.created_at = created_at
    this.updated_at = updated_at
  }
}

export class Recording {
  id
  date
  duration
  channel_name
  path
  constructor(id: string, date: Date, duration: string, channel_name: string, path: string) {
    this.id = id
    this.date = date
    this.duration = duration
    this.channel_name = channel_name
    this.path = path
  }
}
