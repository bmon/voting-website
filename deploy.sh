#!/bin/sh

PROJECT_ID=brendanr-misc
SERVICE_NAME=emote-voting
REGION=asia-northeast1
# load env config from .env file and format it ready for use with gcloud. PORT is a reserved var and must be filtered.
ENVCONF=$(cat .env | grep -v PORT | grep -v GOOGLE_APPLICATION_CREDENTIALS | tr '\n' ',')

gcloud builds submit --tag gcr.io/$PROJECT_ID/$SERVICE_NAME
gcloud run deploy $SERVICE_NAME --image gcr.io/$PROJECT_ID/$SERVICE_NAME --platform managed --set-env-vars="$ENVCONF" --region=$REGION
