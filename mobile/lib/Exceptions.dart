import 'dart:io';

class NotFoundException extends HttpException implements Exception {
  NotFoundException(super.message);
}

class BadRequestException extends HttpException implements Exception {
  BadRequestException(super.message);
}

class UnauthorizedException extends HttpException implements Exception {
  UnauthorizedException(super.message);
}

class ForbiddenException extends HttpException implements Exception {
  ForbiddenException(super.message);
}

class InternalServerErrorException extends HttpException implements Exception {
  InternalServerErrorException(super.message);
}

class MethodNotAllowedException extends HttpException implements Exception {
  MethodNotAllowedException(super.message);
}

class NotAcceptableException extends HttpException implements Exception {
  NotAcceptableException(super.message);
}

class RequestTimeoutException extends HttpException implements Exception {
  RequestTimeoutException(super.message);
}

class TooManyRequestsException extends HttpException implements Exception {
  TooManyRequestsException(super.message);
}

class UnsupportedMediaTypeException extends HttpException implements Exception {
  UnsupportedMediaTypeException(super.message);
}

class PreconditionFailedException extends HttpException implements Exception {
  PreconditionFailedException(super.message);
}

class PayloadTooLargeException extends HttpException implements Exception {
  PayloadTooLargeException(super.message);
}

class UriTooLongException extends HttpException implements Exception {
  UriTooLongException(super.message);
}

class RangeNotSatisfiableException extends HttpException implements Exception {
  RangeNotSatisfiableException(super.message);
}

class ExpectationFailedException extends HttpException implements Exception {
  ExpectationFailedException(super.message);
}

class EnhanceRequestException extends HttpException implements Exception {
  EnhanceRequestException(super.message);
}

class UnavailableForLegalReasonsException extends HttpException implements Exception {
  UnavailableForLegalReasonsException(super.message);
}

class InsufficientStorageException extends HttpException implements Exception {
  InsufficientStorageException(super.message);
}

class ServiceUnavailableException extends HttpException implements Exception {
  ServiceUnavailableException(super.message);
}

class GatewayTimeoutException extends HttpException implements Exception {
  GatewayTimeoutException(super.message);
}

class HttpVersionNotSupportedException extends HttpException implements Exception {
  HttpVersionNotSupportedException(super.message);
}

class VariantAlsoNegotiatesException extends HttpException implements Exception {
  VariantAlsoNegotiatesException(super.message);
}

class TooManyConnectionsException extends HttpException implements Exception {
  TooManyConnectionsException(super.message);
}