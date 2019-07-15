/*
 * {{COPYRIGHT}}
 *
 * Generated on: {{GENERATED_DATE}}
 * {{EAD_COMMENT_NOTICE}}
 */

#include <string.h>
#include <stdio.h>
#include <stdlib.h>

#include "ead_collection.h"
#include "ead_helpers.h"

/* When debugging EAD helpers, then provide EAD_LOG implementation */
#ifndef EAD_LOG
#define EAD_LOG(debug_level, format_message)
#endif

#define EAD_LOG_LEVEL_WARNING 2
#define EAD_LOG_LEVEL_ERROR   3


bool ead_find_entry(const char *location, ead_item_t *item)
{
    if ('/' == *location)
    {
        /* Do ignore the / for paths starting with it */
        location++;
    }

    for (int file = 0; file < STATIC_DATA_INDEX_SIZE; ++file)
    {
        if (strcmp(location, static_data_index[file].location) == 0)
        {
            memcpy(item, &static_data_index[file], sizeof(ead_item_t));
            return true;
        }

    }
    return false;
}


/* Functions imitating file stream functionality */


ead_stream_t* ead_open(const char *filename, const char *mode)
{
    /* Mode is fully ignored, just added here to have fopen signature */
    ead_stream_t *ret = malloc(sizeof(ead_stream_t));
    ead_item_t   item;

    if (!ead_find_entry(filename, &item))
    {
        EAD_LOG(EAD_LOG_LEVEL_ERROR, ("File not found"));
        return NULL;
    }

    memcpy(&ret->content, &item, sizeof(ead_item_t));
    ret->current_location = 0;
    ret->error            = 0;
    return ret;
}


int ead_close(ead_stream_t *stream)
{
    free(stream);
    return 0;
}


int ead_fseek(ead_stream_t *stream, int offset, int origin)
{
    if (SEEK_SET == origin)
    {
        if (offset >= stream->content.size)
        {
            stream->error = 1;
            EAD_LOG(EAD_LOG_LEVEL_WARNING, ("Wrong offset"));
            return -1;
        }
        stream->current_location = offset;

    }
    else if (SEEK_CUR == origin)
    {
        if ( (offset + stream->content.size) >= stream->content.size)
        {
            stream->error = 1;
            EAD_LOG(EAD_LOG_LEVEL_WARNING, ("Wrong offset"));
            return -1;
        }
        stream->current_location = offset + stream->content.size;

    }
    else if (SEEK_END == origin)
    {
        if ( (offset + stream->content.size) < 0 && (offset > 0))
        {
            stream->error = 1;
            EAD_LOG(EAD_LOG_LEVEL_WARNING, ("Wrong offset"));
            return -1;
        }
        stream->current_location = offset + stream->content.size;

    }
    else
    {
        /* wrong origin given */
        stream->error = 1;
        EAD_LOG(EAD_LOG_LEVEL_WARNING, ("Using wrong origin"));
        return -1;
    }

    /* Didn't failed on any error, therefore pass as success */
    return 0;
}


int ead_error(ead_stream_t *stream)
{
    return stream->error;
}


size_t ead_read(void *ptr, size_t size, size_t nmemb, ead_stream_t *stream)
{
    int read_bytes = nmemb * size;

    if (0 >= read_bytes)
    {
        EAD_LOG(EAD_LOG_LEVEL_ERROR, ("Negative 'size' or 'nmemb'"));
        return 0;
    }

    if (read_bytes + stream->current_location >= stream->content.size)
    {
        /* Read less than requested because the "file" is smaller */
        read_bytes = stream ->content.size - stream->current_location;
        EAD_LOG(EAD_LOG_LEVEL_WARNING, ("Requesting to read more than file size"));
    }

    memcpy(ptr, stream->content.data, read_bytes);

    return read_bytes;
}

