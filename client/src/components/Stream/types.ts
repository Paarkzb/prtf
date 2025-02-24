export class Channel {
  id
  rf_user_id
  channel_name
  live
  rf_active_stream_id
  created_at
  updated_at
  constructor(
    id: string,
    rf_user_id: string,
    channel_name: string,
    live: boolean,
    rf_active_stream_id: string,
    created_at: Date,
    updated_at: Date
  ) {
    this.id = id
    this.rf_user_id = rf_user_id
    this.channel_name = channel_name
    this.live = live
    this.rf_active_stream_id = rf_active_stream_id
    this.created_at = created_at
    this.updated_at = updated_at
  }
}
