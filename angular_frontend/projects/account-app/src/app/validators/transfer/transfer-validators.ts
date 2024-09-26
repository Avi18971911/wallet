import { AbstractControl, ValidationErrors, ValidatorFn } from '@angular/forms';

export var exceedBalanceErrorMessage = "amountExceedsBalance";
export function amountLessThanOrEqualToBalance(balance: string): ValidatorFn {
  return (control: AbstractControl): ValidationErrors | null => {
    const amount = control.value;
    if (amount !== null && amount !== undefined && parseFloat(amount) > parseFloat(balance)) {
      return { amountExceedsBalance: true };
    }
    return null;
  };
}
