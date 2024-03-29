/*
 * {{COPYRIGHT}}
 *
 * Generated on: {{GENERATED_DATE}}
 * {{EAD_COMMENT_NOTICE}}
 */

#ifndef EMBEDDED_AUTOGENERATED_DATA_EAD_HELPERS_H_
#define EMBEDDED_AUTOGENERATED_DATA_EAD_HELPERS_H_

#include <stdbool.h>

#include "{{INCLUDE_PREFIX}}/ead_structures.h"

bool ead_find_entry(const char *location, ead_item_t *item);

ead_stream_t* ead_open(const char *filename, const char *mode);

int ead_close(ead_stream_t *stream);

int ead_fseek(ead_stream_t *stream, int offset, int origin);

int ead_error(ead_stream_t *stream);

size_t ead_read(void *ptr, size_t size, size_t nmemb, ead_stream_t *stream);


#endif /* EMBEDDED_AUTOGENERATED_DATA_EAD_HELPERS_H_ */
