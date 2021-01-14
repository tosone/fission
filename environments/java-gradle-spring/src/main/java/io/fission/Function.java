package io.fission;

import org.springframework.http.RequestEntity;
import org.springframework.http.ResponseEntity;

public interface Function<RequestEntity, ResponseEntity> {
	org.springframework.http.ResponseEntity call(org.springframework.http.RequestEntity req, Context context);
}
